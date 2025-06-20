// Map Clustering Performance Tests
const clusteringTests = {
  // Test clustering logic
  testClusteringLogic: () => {
    console.log('TEST: Clustering should group nearby flights');
    
    // Mock flights close together
    const flights = [
      { lat: 40.7128, lng: -74.0060, callsign: 'NYC1' },
      { lat: 40.7589, lng: -73.9851, callsign: 'NYC2' }, // Close to NYC1
      { lat: 51.5074, lng: -0.1278, callsign: 'LON1' },  // London - far
    ];
    
    // Simple clustering logic test
    const clusterDistance = 0.1; // degrees
    const clusters = [];
    
    flights.forEach(flight => {
      const nearbyCluster = clusters.find(cluster => 
        Math.abs(cluster.lat - flight.lat) < clusterDistance &&
        Math.abs(cluster.lng - flight.lng) < clusterDistance
      );
      
      if (nearbyCluster) {
        nearbyCluster.count++;
        nearbyCluster.flights.push(flight);
      } else {
        clusters.push({
          lat: flight.lat,
          lng: flight.lng,
          count: 1,
          flights: [flight]
        });
      }
    });
    
    const hasNYCCluster = clusters.some(c => c.count === 2);
    const hasLondonCluster = clusters.some(c => c.flights[0].callsign === 'LON1');
    
    console.log(`  ✓ NYC flights clustered: ${hasNYCCluster}`);
    console.log(`  ✓ London flight separate: ${hasLondonCluster}`);
    console.log(`  ✓ Total clusters: ${clusters.length}`);
    
    return hasNYCCluster && hasLondonCluster && clusters.length === 2;
  },

  runTests: () => {
    console.log('=== CLUSTERING TESTS ===');
    const result = clusteringTests.testClusteringLogic();
    console.log(`Result: ${result ? '✅ PASS' : '❌ FAIL'}`);
    return result;
  }
};

module.exports = clusteringTests;

if (require.main === module) {
  clusteringTests.runTests();
}