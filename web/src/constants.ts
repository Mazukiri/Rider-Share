export const API_URL = (process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8081').replace(/\/$/, '');
export const WEBSOCKET_URL = (process.env.NEXT_PUBLIC_WEBSOCKET_URL ?? 'ws://localhost:8081/ws').replace(/\/$/, '');
