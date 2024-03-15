# Chain Alert (Work in Progress)

## Description:
### Chain Alert is a blockchain monitoring application that tracks transactions for specified addresses and sends real-time notifications when activity is detected. It aims to provide users with a convenient way to stay informed about their blockchain assets.

## Features (Planned):
- Multi-chain support: Monitor transactions across various blockchain networks (e.g., Ethereum, Polygon ,Binance Smart Chain).
- Configurable notifications: Choose your preferred notification method (email, Slack, Telegram, etc.).
- User interface: Manage monitored addresses and notification preferences through a user-friendly interface (future implementation).
- Persistent data storage: Securely store client configurations and transaction history for reliable access.

## Current Progress (v1):
- ✅ Implemented ETHClient for interacting with the Ethereum network.
- ⏳ Configuring Redis for caching data (in progress).
- Creating a sample notification using Slack integration.
- User API for adding and managing monitored addresses (development in progress).

## Next Steps (v1): 
- Configure Postgress database for storing persistent client configurations.
- Configure Apache Kafka for efficient message queuing and event streaming.
- Develop a Notification Service for handling various notification channels (email, Telegram, etc.).

## Technology Stack:
- Programming Language: Golang and Python
- Blockchain Client Library: go-ethereum
- Caching: Redis
- Database: PostgreSQL
- Message Queue: Apache Kafka
- Notification Channels: Slack (example), Email (to be implemented), Telegram (to be implemented)

