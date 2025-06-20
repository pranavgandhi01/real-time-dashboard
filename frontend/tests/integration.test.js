// Integration Tests - Production Grade
const mapTests = require('./map.test.js');
const markerTests = require('./marker-icons.test.js');

const integrationTests = {
  // Test complete map functionality
  testMapIntegration: () => {
    console.log('=== INTEGRATION TESTS ===');
    
    // Run component tests
    console.log('\n1. Running Map Component Tests...');
    const mapResult = mapTests.runAllTests();
    
    // Run marker icon tests  
    console.log('\n2. Running Marker Icon Tests...');
    const markerResult = markerTests.runAllTests();
    
    // Overall result
    const allPassed = mapResult && markerResult;
    
    console.log('\n=== INTEGRATION RESULTS ===');
    console.log(`Map Component: ${mapResult ? 'PASS' : 'FAIL'}`);
    console.log(`Marker Icons: ${markerResult ? 'PASS' : 'FAIL'}`);
    console.log(`Overall: ${allPassed ? '✅ READY FOR DEPLOYMENT' : '❌ NEEDS FIXES'}`);
    
    return allPassed;
  }
};

module.exports = integrationTests;

if (require.main === module) {
  integrationTests.testMapIntegration();
}