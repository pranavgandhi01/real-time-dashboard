// Frontend Performance Configuration
const performanceConfig = {
  // Map display limits
  maxDisplayFlights: parseInt(process.env.NEXT_PUBLIC_MAX_DISPLAY_FLIGHTS) || 100,
  
  // Clustering configuration
  clusterDistance: parseFloat(process.env.NEXT_PUBLIC_CLUSTER_DISTANCE) || 0.1,
  
  // WebSocket configuration
  reconnectDelay: parseInt(process.env.NEXT_PUBLIC_RECONNECT_DELAY) || 3000,
  maxReconnectAttempts: parseInt(process.env.NEXT_PUBLIC_MAX_RECONNECT_ATTEMPTS) || 5,
  
  // Performance thresholds
  performanceThresholds: {
    flightsPerSecond: 1000,
    messagesPerSecond: 5000,
    maxProcessingTime: 100 // milliseconds
  }
};

export default performanceConfig;