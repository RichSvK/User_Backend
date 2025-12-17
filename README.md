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
- Text Editor: Visual Studio Code

## API Endpoints
### Watchlist Management
- `GET /api/auth/watchlist` - Retrieve user's stock watchlist
- `POST /api/auth/watchlist` - Add stock to watchlist
- `DELETE /api/auth/watchlist/:symbol` - Remove stock from watchlist

### Authentication
- `POST /api/user/register` - Create new user account
- `POST /api/user/login` - User login
- `POST /api/auth/user/logout` - User logout

### User Account
- `GET /api/auth/user/profile` - Get user profile
- `DELETE /api/auth/user/delete` - Delete user account
