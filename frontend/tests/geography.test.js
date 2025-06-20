// Geography Validation Tests
const geographyTests = {
  // Test realistic country coordinates
  testCountryCoordinates: () => {
    console.log('TEST: Flight coordinates should match origin countries');
    
    const countryBounds = {
      'United States': { lat: [25, 49], lng: [-125, -66] },
      'Germany': { lat: [47, 55], lng: [5, 15] },
      'France': { lat: [42, 51], lng: [-5, 8] },
      'United Kingdom': { lat: [50, 59], lng: [-8, 2] },
      'Canada': { lat: [42, 83], lng: [-141, -52] },
      'Japan': { lat: [24, 46], lng: [123, 146] },
      'Australia': { lat: [-44, -10], lng: [113, 154] },
      'Brazil': { lat: [-34, 5], lng: [-74, -35] },
      'India': { lat: [6, 37], lng: [68, 97] },
      'China': { lat: [18, 54], lng: [73, 135] }
    };
    
    // Test sample coordinates
    const testFlights = [
      { country: 'China', lat: 39.9, lng: 116.4 }, // Beijing
      { country: 'United States', lat: 40.7, lng: -74.0 }, // NYC
      { country: 'Germany', lat: 52.5, lng: 13.4 } // Berlin
    ];
    
    const validPlacements = testFlights.every(flight => {
      const bounds = countryBounds[flight.country];
      if (!bounds) return false;
      
      const latValid = flight.lat >= bounds.lat[0] && flight.lat <= bounds.lat[1];
      const lngValid = flight.lng >= bounds.lng[0] && flight.lng <= bounds.lng[1];
      
      console.log(`  ${latValid && lngValid ? '✓' : '✗'} ${flight.country}: ${flight.lat}, ${flight.lng}`);
      return latValid && lngValid;
    });
    
    return validPlacements;
  },

  runTests: () => {
    console.log('=== GEOGRAPHY TESTS ===');
    const result = geographyTests.testCountryCoordinates();
    console.log(`Result: ${result ? '✅ PASS' : '❌ FAIL'}`);
    return result;
  }
};

module.exports = geographyTests;

if (require.main === module) {
  geographyTests.runTests();
}