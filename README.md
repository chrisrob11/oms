How to setup the environment

1) Ensure docker is installed
2) Install golang on your box 
3) In the repo on your box run `make build`
   This will create two binaries 
      - oms_server - Contains a web server that can receive and handle requests
      - omsclient - A executable client that can be used to do the required operations
4) Run `make tools`
   NOTE: this was run on a mac
5) run `docker-compose -up -d`
   NOTE: this will launch postgres locally to be used by the running service
6) run `make migrate-up`
   Under the covers this just runs sql-migrate. It creates the schemas and tables required
7) run `./bin/oms_server`
   This will launch the server locally and connect to postgres
8) To bootstrap the seed data into the system run the following command
    `./bin/omsclient import --file ./placements_teaser_data.json --generateInvoices`
    This step takes 30 seconds to initialize all the data.
    It will output that 419 campaigns, 10000 campaignLineItems and 419 invoices


Bucket 1 - How to access functionality

1) Data from step 1 should be in. There are three tables, Campaigns, CampaignLineItems and Invoices. 
   I made the implementation choice to seed "via the api". So I constructed all of the required operations up front to do so.
2) Each Table has an associated function in the client to do options outlined. The command line has help as well
   - List Campaigns 
      run `./bin/omsclient lc` - Writes all list view of campaigns
   - Show Campaign
      run `./bin/omsclient sc -id 1` - Shows one selected campaign
   - List CampaignLineItems
      run `./bin/omsclient lcli --followNextPage` - Will write out all values
      run `./bin/omsclient lcli` - Will write out the first 500 values (All list commands support paging)
   - Show CampaignLineItem
      run `./bin/omsclient scli -id 1` - Shows one selected campaign line item
   - List Invoices
      run `./bin/omsclient li` - Writes all list view of invoices
   - Show Invoice
      run `./bin/omsclient si -id 1` - Shows one selected invoice line item
   - Adjust Invoice
      run `./bin/omsclient adjust-invoice -id 2 --totalAdjustments 445.00` - Adjusts the specified invoice to the new adjustment

Bucket 2

1) A new Invoice can be generated for any campaign at anytime
   - run `/bin/omsclient generate-invoice -id 200`
   Relationship is 1 campaign can have many invoices
2) Create operations
   - run `./bin/omsclient create-campaign -name "Campaign1"`
   - run `./bin/omsclient ccli --actual 23.23 --adjustments 36336.25 --booked 6337.33 --campaignId 10"`
3) Update operations
   - run `./bin/omsclient uc -id 100 -name "blah my new name"`
   - run `./bin/omsclient ucli -id 10004 -actual 100.00`
4) Export Invoices
   - run `./bin/omsclient li -allFields`
   Kind of did a cheap version, its tab delimited like the existing one, but have a switch to write all fields. Could be used to export data and import into a spreadsheet
5) Unit tests - Not added, and I'm mentioning this as I'm a bit frustrated about not having them. Most of this was because I have so  much dang boilerplate and I was attempting to "get pieces e2e" at the beginning through client commands I ran out of time on this. 

At least I had linters to find errors. Much of the linter setup and make file was from a personal project, just call


6) Archiving - was going to add this. Basically a background task that simply "moves all of a graph of a campaign" to a file and "pushes" it S3, but didn't get there.


Bucket 3

What I wanted to do is simply make a jsonnet manifest for installing the application onto kubernetes. The idea being that this can be used to install the application on a kubernetes cluster where ever. 

Why this? Well I started using kubernetes in 2014, very early on. I've worked with it alot. I've made "operators" as well. So I wanted to at least throw something together that shows I can do this.

Its not completely working. What I did was create a jsonnet file ./deployments/kubernetes/app.jsonnet

Using a script `./scripts/generate_app_manifest.sh` it will create the yaml output for all the things for the service. In order for the tool to work you will need to install `jsonnet`. Brew or apt has packages to install it.

This can be used to deploy to kubernetes `kubectl create -f ./deployments/kubernetes/app.yaml`

The app doesn't work BECAUSE it needs to run the migration against the postgres database. Ideally this is a task that is setup to just run. I could get this working by port-forwarding to postgres on the local kubernetes cluster. Tweak my migration file and run it. Then things would work as expected. Its close.

Thoughts on the excercise

1) I have not created a completely new CRUD app in golang ever really.
2) Tools that I use on a normal basis aren't the best for that. I wouldn't use sqlc, I'd just do dynamic sql generation in a safe manner. There are a couple other tools I could have used, SqlBoiler, go-jet.
3) Lots of the code is boiler plate. I wish I had more time to write unit tests and dry it up. 
4) Logging (and tracing and metrics) took a back seat. I added basic logger, but should have made its usage more universal.
5) Missed adding a migration piece to run prior to everything.
6) When the data is uploaded there is a precision issue, the data is not an exact match. Its close. I think its just an issue with the precision conversion from float64->string for db insert, then to numeric. Its likely just a bug in my code in the strconv.FormatFloat that I have. I didn't have time to fix this. Only noting this problem up front.

Regardless of the outcome I'm glad I've done this. A good crash course in trying to build something quicker. I would have liked to play with react but decided instead to just make a cli.


