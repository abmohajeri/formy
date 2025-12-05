# Formy

The core engine of Formy - a form submission service that sends form data directly to Telegram. Formy Core provides a
secure, scalable backend for handling form submissions with built-in CAPTCHA protection, rate limiting, and domain
whitelisting.

## Features

- **Form Submission API**: Simple REST API endpoint for submitting form data
- **Telegram Integration**: Automatic delivery of form submissions to Telegram chats
- **CAPTCHA Protection**: Support ALTCHA to prevent spam
- **Rate Limiting**: Built-in rate limiting (30 requests per minute per IP) to prevent abuse
- **Domain Whitelisting**: Restrict form submissions to allowed domains only
- **Form Tokens**: Unique tokens for each form with UUID-based identification
- **CORS Support**: Cross-origin resource sharing enabled for web forms
- **Redirect Support**: Custom redirect URLs after successful form submission
- **CC Support**: Send form submissions to multiple Telegram chats

## Prerequisites

- **Go**: Version 1.22.5 or higher
- **PostgreSQL**: Database server (version 12 or higher recommended)
- **Telegram Bot**: A Telegram bot token from [@BotFather](https://t.me/BotFather)
- **ALTCHA**: ALTCHA HMAC key for CAPTCHA solution

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd core
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**

   Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

   Edit `.env` and fill in your configuration values. See [Configuration](#configuration) section for details.

4. **Set up PostgreSQL database**

   Create a PostgreSQL database:
   ```sql
   CREATE DATABASE formy_db;
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8030` by default.

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and update the following
variables:

### Basic Configuration

- `BASE_URL`: Base URL of your application

### Database Configuration

- `DATABASE_HOST`: PostgreSQL host (default: `localhost`)
- `DATABASE_USERNAME`: Database username
- `DATABASE_PASSWORD`: Database password
- `DATABASE_DB_NAME`: Database name (default: `formy_db`)
- `DATABASE_PORT`: Database port (default: `5432`)

### Telegram Bot Configuration

- `TELEGRAM_PROXY_URL`: URL or proxy URL for Telegram API used as base url of telegram client
- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token from [@BotFather](https://t.me/BotFather)
- `TELEGRAM_DEBUG`: Enable debug mode (`true`/`false`)

### CAPTCHA Configuration

- `ALTCHA_HMAC_KEY`: ALTCHA HMAC key for challenge generation and verification

## Usage

### API Endpoints

#### Submit Form Data

```
POST /{FORM_TOKEN}
```

Submit form data using a form token UUID.

**Example Request:**

```html

<form action="BASE_URL/FORM_TOKEN" method="post">
    <input type="email" name="email">
    <textarea name="message"></textarea>
    <input type="submit" value="Send">
</form>
```

#### Get CAPTCHA Challenge

```
GET /captcha
```

Returns an ALTCHA challenge for client-side CAPTCHA verification.

#### Telegram Webhook

```
POST /{telegram_bot_token}
```

Webhook endpoint for Telegram bot updates. Automatically configured on startup.

### Form Token Management

Form tokens are managed through the Telegram bot interface. Each token:

- Has a unique UUID
- Is associated with a user and Telegram chat
- Can have multiple allowed domains
- Supports CC functionality to forward submissions to other chats

### Domain Whitelisting

Each user can configure allowed domains for their form tokens. Only submissions from whitelisted domains will be
accepted.

### Rate Limiting

The API implements rate limiting:

- **Limit**: 30 requests per minute per IP address
- **Response**: HTTP 429 (Too Many Requests) when limit is exceeded

## Project Structure

```
core/
├── config/              # Configuration and middleware
│   ├── middleware.go   # Rate limiting middleware
│   └── postgres.go     # Database connection setup
├── controllers/        # Request handlers
│   ├── api/           # API controllers
│   └── telegram/      # Telegram webhook handlers
├── models/            # Database models
│   ├── User.go        # User model
│   ├── FormToken.go   # Form token model
│   └── AllowedDomain.go # Allowed domain model
├── routes/            # Route definitions
│   └── router.go      # Main router setup
├── services/          # Business logic services
│   ├── TelegramService.go  # Telegram bot service
│   ├── CaptchaService.go   # CAPTCHA verification
│   ├── FormTokenService.go # Form token management
│   └── DomainService.go    # Domain management
├── utils/             # Utility functions
│   ├── RequestUtils.go    # Request helpers
│   └── UUIDUtils.go       # UUID utilities
├── views/             # HTML templates
│   ├── form-template.html      # Form template
│   ├── form-verification.html  # Success/error page
│   └── assets/        # Static assets
├── main.go            # Application entry point
├── go.mod             # Go module dependencies
└── .env.example       # Environment variables template
```

## Development

### Running in Development Mode

```bash
go run main.go
```

### Building

```bash
go build -o formy-core main.go
```

### Database Migrations

Database migrations are handled automatically using GORM AutoMigrate. The following models are migrated on startup:

- `User`
- `FormToken`
- `AllowedDomain`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. When contributing:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows Go best practices and includes appropriate tests where applicable.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, questions, or contributions, please open an issue on the GitHub repository.
