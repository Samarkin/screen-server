# ScreenServer API Description

#### `GET /api/health`
Get information about the server.

#### `GET /api/messages`
Get screen contents.

#### `POST /api/messages`
Display message on the next line.

#### `DELETE /api/messages`
Clear entire screen.

#### `GET /api/messages/{line}`
Get message displayed on the given line.

#### `PUT /api/messages/{line}`
Display message on the given line.

#### `DELETE /api/messages/{line}`
Clear the given line.
