export const API = {
  RIDE: (id: string) => `/api/rides/${id}`,
  RIDES: '/api/rides',
  LOCATIONS: '/api/locations',
  USERS: {
    MAIN: '/api/users',
    CONTACT: '/api/users/contact',
  },
  AUTH: {
    LOGIN: '/api/auth/login',
    VERIFY: '/api/auth/verify',
    USER: '/api/auth/user',
    LOGOUT: '/api/auth/logout',
  },
}
