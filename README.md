# Introduction 
The purpose of this document is to describe the subscrintion manager module, its functions, and design.

# Overview
The pricing manager is a backend service that implements billing, payment history, and subscription features. It implements them through integration with the payment gateway stripe.com

# Design
From the user perspective, the module will have the following pages:
- **Pricing**. The page is available to both registered and unauthorized users and describes all available plans.
- **Plans**. The page is available to registered users only (as a part of user account?) and allows the user to upgrade/downgrade/cancel their current plan.
- **Billing**. The page is available to registered users only (as a part of user account?) and allows the user to view upcoming/current billing statement as well as the history of payments.

To accept payments the module will use [Stripe Checkout](https://stripe.com/docs/payments/checkout) integration method (from the [3 available](https://stripe.com/docs/payments/online-payments)).

## Database Structure
The service has the following table structure:
```plantuml
object transactions {
  id: uuid
  userId: uuid
  priceId: varchar(100)
  userPlanId: uuid
  sessionId: varchar(100)
  customerId: varchar(100)
  status: varchar(50)
  created_ts: timestamp
  last_modified_ts: timestamp
}

object user_plans {
  id: uuid
  userId: uuid
  planId: int4
  customerId: varchar(100)
  priceId: varchar(100)  
  subscriptionId: varchar(100)
  status: varchar(50)
  created_ts: timestamp
  last_modified_ts: timestamp
}

object available_plans {
  id: int4
  name: varchar(50)
  description: varchar(250)
  price: numeric
  recurrence: int4
  priceId: varchar(100)
}


transactions "N " --> "1 " user_plans
user_plans "N " --> "1 " available_plans
```

## Pricing Component
The pricing page on the web-site as well the portal component pricing.tsx should be implemented as per this [design](https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490).

## Plans Component - Frontend
**What needs to be done**
1. In officekube portal repo modify the code of the pricing.tsx component to handle the logic of a user switching from one plan to another as follows:
   - The component should disable the button Sign Up button for the plan that the user is currently on. To determine what the current plan is the component should call an API endpoint of the subscription manager service GET /plans/current at the time of its loading.
   - Replace button Notify Me for the plan Solo with the button SIGN UP.
   - When the user clicks on Sign Up button of any other plan a Switch Plan dialog should pop up (refer to this [design](https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490)). The dialog should be created using [Syncfusion React library](https://ej2.syncfusion.com/react/documentation/dialog/getting-started/).
   - In the dialog when the user clicks button Switch make a GET call to the endpoint /payments/checkout of the subscription manager backend (see below) and pass a query parameter named price_id with the value that depends on which button SIGN UP has been clicked by the user before they got to the dialog. For the plan Enthusiast that value should be set "free", for the plan Solo it should be "price_1LECZyKUSkDFrC1EroX3h7NW".
   - If the user clicked Cancel simply return them back to the Plans page.

### Success and Failure Pages
After a user has been redirected to the Stripe checkout page, Stripe will redirect the user back to either a success or a failure page indicating whether the user has successfully signed up for our subscription.
The pricing.tsx will be responsible for showing either success or failure. For that to work the page might receive a URL parameter named checkoutResult. Modify the page as follows:
1. Use the react-router-dom library to retrieve the value of the URL parameter checkoutResult immediately after the page has been loaded into a web-browser.
2. If the parameter is equal to "success" then show a popup message (using Dialog component from Syncfusion library) with a button OK and a message "Thank you for your subscription!". When the user clicks OK, the dialog should be closed.
3. If the parameter is equal to "failure" then show a popup message (using Dialog component from Syncfusion library) with a button OK and a message "Sorry, something went wrong. Please try again later or contact us!". When the user clicks OK, the dialog should be closed.
4. If the parameter is not set to any value then no action should be taken.

## Plans Component - Backend
**What needs to be done**
The service is based on the following stack:
- Language: go-lang
- Web Framework: Gin Web Framework
- Configuration Manager: Viper
- ORM: GORM
- DB Backend: PostgreSQL
- Code Generator: Open API Code Generator

### 1. Endpoint GET /plans
As a part of this assignment implement /plans  (refer [openapi.yml](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/blob/main/api/openapi.yml)). Assume that all available plans are stored in a table AvailablePlans persisted in the PostgreSQL db to which the service has read/write access.
The service should make a call into the DB and retrieve all records from the mentioned table where field active is equal to "true". The table AbailablePlans has the following structure:

- id int primary key auto incremental
- name char 50 required
- description char 250
- price real required
- recurrence int required default 30

For testing purpose the table can be populated with the following records:

|id|name|description|price|recurrence|
|-|-|-|-|-|
|1|Enthusiast||0|30|
|2|Solo||10|30|
|3|Expert||30|30|
|4|Team||100|30|

**Development Approach**
1. Generate a code skeleton for the application using the following command:
openapi-generator-cli generate --package-name workspaceEngine -g go-gin-server -i openapi.yml 
1. Implement the endpoint GET /plans using Gin Web Framework and GORM ORM (for access to DB).
1. Avoid hard-coding service configuration (e.g. db connection parameters). The service configuration should be persisted in a YAML file subscription_manager.yml.


### 2. Endpoint GET /payments/checkout
The endpoint creates new subscription (for new users), upgrades/downgrades the subscription, cancels it (in case a user selects a free plan).

The following activity diagram provides details around the logic implemented in the endpoint.

```plantuml
start
:Retrieving the stripe key;
:Get user id (GetUserId());
:Pull record from user_plans 
 with user_id;
If (Record exists?) then (YES)
  If (record.priceId == price_id?) then (YES)
    :Return HTTP 208;
    end
    Else (NO - Switch to Another Plan)
      :Pull record from available_plans 
          (priceId = price_id);
      If (price_id == priceId of 
        available plan 
        with price $0?) then (NO)
      If (user_plans.subscriptionId is empty?) then (NO)
        :Upgrade/Downgrade 
             Subscription;
      Else (Upgrade from free to paid plan)
        :Create checkout session with the customerId;
        :Create new record in transactions;
        :Update user_plans record 
     (priceId, planId, last_modified_ts);
        :Return HTTP 302;
        end
      Endif
    Else (YES)
      :Cancel Current 
        Subscription;
      :Set user_plans.subscriptionId to null;
    EndIf
    :Create new record in transactions;
    :Update user_plans record 
     (priceId, planId, last_modified_ts);
    :Return HTTP 200;
    end
  Endif
Else (NO - New Subscription)
  :Call "Create a customer" API;
  :Create new record in user_plans 
  (with customerId);
  If (price_id == priceId of 
      available plan 
      with price $0?) then (NO)
    :Create checkout session with the customerId;
    :Create new record in transactions;
    :Return HTTP 302;
    end
  Else (YES)
    :Create new record in transactions
      without sessionId;
    :Set status to CURRENT for 
    user_plans & transactions records;
    :Return HTTP 200;
    end
  EndIf
EndIf
```


### 3. Endpoint POST /payments/stripewebhook
The endpoint is called by the Stripe backend as a part of the checkout process. Refer to the [stripe guide](https://stripe.com/docs/payments/checkout/fulfill-orders) for the Stripe Checkout integration.

### 4. Endpoint GET /plans/current
The endpoint /payments/current is implemented as follows:
- Secure the endpoint with a call to IsApiAuthenticated().
- Pull a user id using the function GetUserId
- Retrieve a record from user_plans where userId == user id and status == 'CURRENT'.
- If no record is found then return an http code 404.
- If a record has been found using its field planId retrieve a matching record from the table available_plans.
- Create an instance of the model APlan, populate its properties with proper values from the user_plans and available_plans records and return the model along with http code 200.


### 5. Cancelling Subscription
At this point we won't build subscription cancellation explicitly for 2 reasons:
1. A user would have an option to switch back to the free plan.
2. If a user insists on cancellation (i.e. closing the account) we will use a [manual option](https://stripe.com/docs/billing/subscriptions/cancel) to cancel their subscription.

## Billing - Frontend
1. In officekube portal in settings, under the menu item Plan (but above the item Notifications) add a new item called Billing.
2. Similar to the component Pricing.tsx add a component Billing.tsx (in folder src\app\components\Account) that should be designed as per the design frame [Payment History](https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=0%3A1).
3. The Billing.tsx component should invoke the endpoint GET /payments/history. Refer to the [openapi.yml](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/blob/main/api/openapi.yml) spec for details on how to call the endpoint and process its response.

## Billing - Backend
### 1. Endpoint /payments/history
Implement the endpoint /payments/history as follows:
- Create the endpoint handler in a separate file go/api_payments_history.go and using the GIN web framework. 
- Secure the endpoint with a call to IsApiAuthenticated().
- Pull a user id using the function GetUserId
- Retrieve a record from user_plans where userId == user id and status == 'CURRENT'.
- If no record is found then return an http code 404.
- If a record has been found using its field customerId retrieve a list of payments from Stripe using its API endpoint [PaymentIntents](https://stripe.com/docs/api/payment_intents/list).
- Create an array APayment instances return it (as json payload { "payments": array }) along with http code 200.
