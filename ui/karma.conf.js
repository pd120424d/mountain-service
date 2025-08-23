module.exports = function (config) {
    config.set({
      basePath: '',
      frameworks: ['jasmine', '@angular-devkit/build-angular'],
      plugins: [
        require('karma-jasmine'),
        require('karma-chrome-launcher'),
        require('karma-jasmine-html-reporter'),
        require('karma-coverage'),
        require('@angular-devkit/build-angular/plugins/karma')
      ],
      client: {
        jasmine: {},
        clearContext: false // leave Jasmine Spec Runner output visible in browser
      },
      coverageReporter: {
        dir: require('path').join(__dirname, './coverage/ui'),
        subdir: '.',
        reporters: [
          { type: 'text-summary' },
          { type: 'json-summary' } // important for test-cover-ui.sh
        ],
        fixWebpackSourcePaths: true
      },
      reporters: ['progress', 'kjhtml', 'coverage'],
      port: 9876,
      colors: true,
      logLevel: config.LOG_INFO,
      autoWatch: true,
      browsers: ['ChromeHeadlessNoSandbox'],
      customLaunchers: {
        ChromeHeadlessNoSandbox: {
          base: 'ChromeHeadless',
          flags: [
            '--no-sandbox',
            '--disable-setuid-sandbox',
            '--disable-dev-shm-usage',
            '--disable-gpu',
            '--remote-debugging-port=9222'
          ]
        }
      },

      singleRun: false,
      restartOnFileChange: true
    });
  };
  