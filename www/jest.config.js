module.exports = {
  preset: 'ts-jest',
  moduleFileExtensions: [
    'js',
    'ts',
    'json',
    'vue'
  ],
  transform: {
    '^.+\\.ts$': 'ts-jest',
    "^.+\\.js$": "babel-jest",
    '^.+\\.vue$': '@vue/vue3-jest'
  },
  testEnvironment: 'jsdom',
  setupFiles: ['./tests/setup.js'],
  testEnvironmentOptions: {
    customExportConditions: ["node", "node-addons"],
  },
}
