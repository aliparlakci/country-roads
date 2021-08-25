export const API = () => {
  const URI = process.env.BACKEND_URI || "http://localhost:8080";
  
  return {
    RIDES: `${URI}/api/rides`,
    LOCATIONS: `${URI}/api/locations`,
  };
};
