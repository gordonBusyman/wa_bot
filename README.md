# GoLang Telegram Bot for Product Rating

This GoLang-based Telegram bot is designed for businesses to gather product ratings and feedback from customers. It utilizes semantic analysis to interpret user feedback, providing valuable insights into customer satisfaction and product performance. Additionally, the bot integrates an HTTP server to handle order creation, and it sends notifications to users via Telegram.

## Features

- **Product Rating**: Users can rate products and provide feedback directly through Telegram.
- **Semantic Analysis**: The bot analyzes user feedback for sentiment, helping businesses understand customer opinions.
- **Order Creation via HTTP**: External systems can create product orders through HTTP requests, and the bot will notify the relevant user on Telegram.
- **User Notifications**: Sends automated messages to users on Telegram regarding their orders and requests for product feedback.

## Dependencies

- **Chi Router**: A lightweight, idiomatic, and composable router for building Go HTTP services.
- **Telegram Bot API (tgbotapi)**: A Go package that wraps the Telegram Bot API. [GitHub Repository](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **PostgreSQL**: A powerful, open-source object-relational database system.
- **Squirrel**: A fluent SQL generator for Go. [GitHub Repository](https://github.com/Masterminds/squirrel)

## Getting Started

### Prerequisites

- GoLang (version 1.21.0 or later)
- PostgreSQL (version 12 or later)
- A Telegram bot token (obtainable through [BotFather](https://t.me/botfather) on Telegram)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/gordonBusyman/wa_bot
   cd wa_bot
   ```

2. **Set Up Environment Variables**

   Create a `.env` file in the root directory and add the following:

   ```env
   CONNECTLY_BOT_TOKEN=your_telegram_bot_token
   CONNECTLY_DB_HOST=your_database_url
   CONNECTLY_DB_USERNAME=your_database_user
   CONNECTLY_DB_PASSWORD=your_database_pass
   ```

3. **Install Dependencies**

   ```bash
   go get -u all
   ```

4. **Initialize Database**

   Set up your PostgreSQL database and run the provided SQL scripts to create the necessary tables.

5. **Run the Bot**

   ```bash
   go run main.go
   ```

### Usage

- **Interacting with the Bot**: Users can start a conversation with the bot on Telegram to rate products and provide feedback.
- **Creating Orders via HTTP**: Send a POST request to `/orders/{user_id}` with the order details. The bot will process the order and send a notification to the user on Telegram.

## Database Schema
![alt text](https://github.com/gordonBusyman/wa_bot/blob/db.png?raw=true)

## Configuration

- **Bot Token**: Ensure your bot token is correctly set in the `.env` file.
- **Database Connection**: Configure your PostgreSQL connection string in the `.env` file.

## Contributing

Contributions to the project are welcome. Please follow the standard fork-and-pull request workflow.

## License

This project is licensed under the [Apache License](LICENSE).

## Contact

For any queries or contributions, please contact the repository owner.
