export const API = () => {
  const URI = process.env.BACKEND_URI || "http://localhost:8080";

  const withHost = (endpoint: string) => `${URI}${endpoint}`;

  return {
    RIDE: (id: string) => withHost(`/api/rides/${id}`),
    RIDES: withHost("/api/rides"),
    LOCATIONS: withHost("/api/locations"),
  };
};
