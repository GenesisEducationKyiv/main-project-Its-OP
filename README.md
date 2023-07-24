# Genesis KMA SE School 3.0
## Prostakov Oleh
&nbsp;
# !!! Important
**The application relies on SendGrid for sending emails. In order to enable sending emails by the application, you need to provide details of your SendGrid sender to the application, using environmental variables.**

**The next data needs to be passed to the application:**
1. Sender's name, via `SENDGRID_API_SENDER_NAME`.
2. Sender's email via `SENDGRID_API_SENDER_EMAIL`.
3. API key via `SENDGRID_API_KEY`.

**Information about your senders can be taken from this link:** https://app.sendgrid.com/settings/sender_auth/senders.
**The API key can be taken from here:** https://app.sendgrid.com/settings/api_keys.

**Placeholders for all these variables have been added to `docker-compose`. If you use it to run the application, please set the real values instead of the placeholders.**
&nbsp;
# Get Started
1. Clone the repository;
2. run `docker compose up` at the root folder of the project;
3. When the compose stack launches, go to http://localhost:8080/swagger/index.html.

If you create the application container manually from the image, and not via `docker-compose`, please don't forget to specify all the environmental variables mentioned. Also, don't forget to map the port exposed by the container - 8080.

# About the project
The application is an API service, which serves 3 purposes:
1. It allows to get the current rate of BTC to UAH;
2. It allows to subscribe an email to a mailing list, in order to get updates about the rate changes;
3. It allows to send an email with the latest rate of BTC to UAH to all the subscribers.

The rate of BTC to UAH is taken from the BTC/UAH traiding pair on Binance. The UI of the application is served by swagger. The subscribers' emails are stored in a file, which is persisted among different application launches (when ran via `docker-compose`). The API of the application was implemented with the **gin** framework, and SendGrid was choosen as an email provider.

## Architecture
The project was designed according to the principles of Clean Architercture. In consists of 4 packages - **domain**, which is responsible for the rules of the system, **application**, which stores constructs, that provide logic of the application, **infrastructure**, that consists of adapters for the 3rd-party systems, and **web**, which serves the API.

Speaking about the dependencies, **domain** has none, **application** depends on **domain**, while **infrastructure** and **web** depend on both **domain** and **application**. Furthermore, all the services depend on Contracts, so the concrete implementations can be replaced with ease. Such approach simplifies refactorings of the system, and provides an easy way for testing of each construct.
