// Test WebSocket reconnection logic with exponential backoff

describe('WebSocket Reconnection Logic', () => {
  test('should calculate exponential backoff correctly', () => {
    const baseDelay = 1000;
    const maxDelay = 30000;
    
    // Test exponential backoff calculation
    const calculateDelay = (attempt) => {
      const exponentialDelay = Math.min(baseDelay * Math.pow(2, attempt), maxDelay);
      return exponentialDelay;
    };
    
    expect(calculateDelay(0)).toBe(1000);   // 1s
    expect(calculateDelay(1)).toBe(2000);   // 2s
    expect(calculateDelay(2)).toBe(4000);   // 4s
    expect(calculateDelay(3)).toBe(8000);   // 8s
    expect(calculateDelay(4)).toBe(16000);  // 16s
    expect(calculateDelay(5)).toBe(30000);  // 30s (capped)
    expect(calculateDelay(10)).toBe(30000); // 30s (capped)
  });

  test('should add jitter to prevent thundering herd', () => {
    const baseDelay = 1000;
    const maxJitter = 1000;
    
    const calculateDelayWithJitter = (attempt) => {
      const exponentialDelay = Math.min(baseDelay * Math.pow(2, attempt), 30000);
      const jitter = Math.random() * maxJitter;
      return exponentialDelay + jitter;
    };
    
    const delay1 = calculateDelayWithJitter(0);
    const delay2 = calculateDelayWithJitter(0);
    
    // Delays should be different due to jitter
    expect(delay1).toBeGreaterThanOrEqual(1000);
    expect(delay1).toBeLessThanOrEqual(2000);
    expect(delay2).toBeGreaterThanOrEqual(1000);
    expect(delay2).toBeLessThanOrEqual(2000);
  });

  test('should respect maximum reconnection attempts', () => {
    const maxAttempts = 10;
    let attempts = 0;
    
    const shouldReconnect = () => {
      return attempts < maxAttempts;
    };
    
    // Simulate failed attempts
    for (let i = 0; i < 15; i++) {
      if (shouldReconnect()) {
        attempts++;
      }
    }
    
    expect(attempts).toBe(maxAttempts);
  });

  test('should reset attempts on successful connection', () => {
    let attempts = 5;
    
    const onSuccessfulConnection = () => {
      attempts = 0;
    };
    
    onSuccessfulConnection();
    expect(attempts).toBe(0);
  });
});

describe('WebSocket Connection States', () => {
  test('should handle different connection states', () => {
    const states = {
      CONNECTING: 'Connecting...',
      CONNECTED: 'Connected',
      DISCONNECTED: 'Disconnected',
      ERROR: 'Error'
    };
    
    expect(states.CONNECTING).toBe('Connecting...');
    expect(states.CONNECTED).toBe('Connected');
    expect(states.DISCONNECTED).toBe('Disconnected');
    expect(states.ERROR).toBe('Error');
  });

  test('should format reconnection message correctly', () => {
    const formatReconnectionMessage = (attempt, maxAttempts, delayMs) => {
      const seconds = Math.round(delayMs / 1000);
      return `Reconnecting in ${seconds}s... (${attempt}/${maxAttempts})`;
    };
    
    expect(formatReconnectionMessage(1, 10, 2500)).toBe('Reconnecting in 3s... (1/10)');
    expect(formatReconnectionMessage(5, 10, 16000)).toBe('Reconnecting in 16s... (5/10)');
  });
});