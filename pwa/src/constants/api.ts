export const API = {
  RIDE: (id: string) => `/api/rides/${id}`,
  RIDES: '/api/rides',
  LOCATIONS: '/api/locations',
  USERS: {
    MAIN: '/api/users',
  },
  CONTACT: (id: string) => `/api/contact/${id}`,
  AUTH: {
    LOGIN: '/api/auth/signin',
    VERIFY: '/api/auth/verify',
    USER: '/api/auth/user',
    LOGOUT: '/api/auth/logout',
  },
}
