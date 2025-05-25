func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user").(*jwt.Token)
        claims := user.Claims.(jwt.MapClaims)

        if claims["role"] != "admin" {
            return echo.ErrUnauthorized
        }

        return next(c)
    }
}
