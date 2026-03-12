## User Backend Service

### Program Description
User backend is a backend systems for user and stock watchlist service.

## Related Repositories
- **iOS Application**: https://github.com/RichSvK/StockBalance
- **API Gateway**: https://github.com/RichSvK/API_Gateway
- **User and Watchlist services**: https://github.com/RichSvK/User_Backend
- **Stock Services**: https://github.com/RichSvK/Stock_Backend

### System Requirements
Software used in developing this program:
- Go
- Fiber Web Framework
- PostgreSQL
- Redis

## API Endpoints
### Watchlist Management
- `GET /api/v1/auth/watchlist` - Retrieve user's stock watchlist
- `POST /api/v1/auth/watchlist` - Add stock to watchlist
- `DELETE /api/v1/auth/watchlist/:stock` - Remove stock from watchlist

### Authentication
- `POST /api/v1/users/register` - Create new user account
- `POST /api/v1/users/login` - User login
- `POST /api/v1/auth/user/logout` - User logout

### User Account
- `GET /api/v1/auth/users/profile` - Get user profile
- `DELETE /api/v1/auth/users` - Delete user account by admin
- `GET /api/v1/users/verify` - Verify user account