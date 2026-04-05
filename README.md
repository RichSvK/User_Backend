## User Backend Service

### Program Description
User backend is a backend systems for user and stock watchlist service.

## Related Repositories
- **iOS Application**: https://github.com/RichSvK/StockBalance
- **Gateway**: https://github.com/RichSvK/API_Gateway
- **User and Watchlist services**: https://github.com/RichSvK/User_Backend
- **Stock Services**: https://github.com/RichSvK/Stock_Backend

### System Requirements
Software used in developing this program:
- Go
- Fiber Web Framework
- PostgreSQL
- Redis

## API Endpoints
### Authentication
- `POST /api/v1/users/register` - Create new user account
- `POST /api/v1/users/login` - User login
- `POST /api/v1/users/logout` - User logout

### User Account
- `GET /api/v1/auth/users/profile` - Get user profile
- `GET /api/v1/auth/verify` - Verify user account
- `DELETE /api/v1/users` - Delete user account by admin

### Watchlist Management
- `GET /api/v1/watchlists` - Retrieve user's stock watchlist
- `POST /api/v1/watchlists/stocks` - Add stock to user watchlist
- `DELETE /api/v1/watchlists/stocks/:stock` - Remove stock from user watchlist

### Favorites
- `GET /api/v1/favorites` - Retrieve user underwriter favorites
- `POST /api/v1/favorites` - Add underwriter to user favorites
- `DELETE /api/v1/favorites/:underwriter` - Remove an underwriter from user favorites