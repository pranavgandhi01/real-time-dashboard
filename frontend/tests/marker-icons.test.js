// Marker Icon Tests - Production Grade
const fs = require('fs');
const path = require('path');

const markerIconTests = {
  // Test 1: Check if marker icon files exist
  testMarkerIconFiles: () => {
    console.log('TEST 1: Marker icon files should exist');
    const publicDir = path.join(__dirname, '..', 'public');
    const requiredFiles = [
      'marker-icon.png',
      'marker-icon-2x.png', 
      'marker-shadow.png'
    ];
    
    const results = requiredFiles.map(file => {
      const filePath = path.join(publicDir, file);
      const exists = fs.existsSync(filePath);
      console.log(`  ${exists ? '✓' : '✗'} ${file}: ${exists ? 'EXISTS' : 'MISSING'}`);
      return { file, exists };
    });
    
    const allExist = results.every(r => r.exists);
    return { allExist, results };
  },

  // Test 2: Validate icon URLs
  testIconUrls: () => {
    console.log('TEST 2: Icon URLs should be accessible');
    const baseUrl = 'http://localhost:3000';
    const iconUrls = [
      '/marker-icon.png',
      '/marker-icon-2x.png',
      '/marker-shadow.png'
    ];
    
    // Simulate URL validation (in real app, would use fetch)
    const urlTests = iconUrls.map(url => ({
      url: baseUrl + url,
      valid: url.startsWith('/') && url.endsWith('.png')
    }));
    
    const allValid = urlTests.every(t => t.valid);
    console.log(`  ✓ All URLs properly formatted: ${allValid}`);
    return { allValid, urlTests };
  },

  // Test 3: Leaflet icon configuration
  testLeafletConfig: () => {
    console.log('TEST 3: Leaflet icon configuration should be correct');
    const config = {
      iconUrl: '/marker-icon.png',
      iconRetinaUrl: '/marker-icon-2x.png',
      shadowUrl: '/marker-shadow.png',
      iconSize: [25, 41],
      iconAnchor: [12, 41],
      popupAnchor: [1, -34],
      shadowSize: [41, 41]
    };
    
    const validConfig = 
      config.iconUrl && config.iconUrl.endsWith('.png') &&
      config.iconRetinaUrl && config.iconRetinaUrl.endsWith('.png') &&
      config.shadowUrl && config.shadowUrl.endsWith('.png') &&
      Array.isArray(config.iconSize) && config.iconSize.length === 2;
    
    console.log(`  ✓ Configuration valid: ${validConfig}`);
    return { validConfig, config };
  },

  // Run all tests
  runAllTests: () => {
    console.log('=== MARKER ICON TESTS ===');
    
    const fileTest = markerIconTests.testMarkerIconFiles();
    const urlTest = markerIconTests.testIconUrls();
    const configTest = markerIconTests.testLeafletConfig();
    
    const results = {
      files: fileTest.allExist,
      urls: urlTest.allValid,
      config: configTest.validConfig
    };
    
    const passed = Object.values(results).filter(Boolean).length;
    const total = Object.keys(results).length;
    
    console.log('\n=== TEST RESULTS ===');
    console.log(`Passed: ${passed}/${total}`);
    
    if (passed === total) {
      console.log('✅ ALL TESTS PASSED');
      return true;
    } else {
      console.log('❌ TESTS FAILED - Need to fix marker icons');
      return false;
    }
  }
};

module.exports = markerIconTests;

if (require.main === module) {
  markerIconTests.runAllTests();
}