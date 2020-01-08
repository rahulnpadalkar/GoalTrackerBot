# GoalTrackerBot
#### Tracking your progress made easy.

Goal Tracker Bot is a telegram bot written in Go to track your goal progress.


#### Steps to get it up and running

1. Creating a Google Cloud project.
    1. Navigate to [Google Cloud Console Dashboard](https://console.cloud.google.com/home/dashboard)
    2. Create a new project and select it from the project selection screen.
    3. Navigate to IAM & admin section and create a service account.
    4. Download the service account credentials in JSON format and place it in the project root.
    
2. Create a telegram bot account by chatting with [BotFather](https://telegram.me/botfather)

3. Create a Google Spreadsheet and share it with the service account created in step 1.

4. Finally create a .env file with below values
    1. PRIVATE_BOT_TOKEN - Token that BotFather returns you when you create a bot account.
    2. SPREADSHEET_ID - The part of the shared sheet url between d/ and /edit
    3. CREDS_LOC - Path to the credentials saved in step 1 substep 4.
    
5. Run the bot using godotenv -f .env go run main.go



  
#### In case of any query please DM me on [twitter](https://twitter.com/rahulnpadalkar). If you like it you can buy me a :coffee:
