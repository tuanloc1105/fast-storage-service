/**
 * This module is responsible on handling all the inter process communications
 * between the frontend to the electron backend.
 */

import { app, ipcMain } from 'electron';
import { environment } from '../../environments/environment';
import * as fs from 'fs';
import * as path from 'path';

export default class ElectronEvents {
  static bootstrapElectronEvents(): Electron.IpcMain {
    return ipcMain;
  }
}

// Retrieve app version
ipcMain.handle('get-app-version', (event) => {
  console.log(`Fetching application version... [v${environment.version}]`);

  return environment.version;
});

// Write logging messages to file
ipcMain.handle('log', (event, message) => {
  const logFile = path.join(app.getPath('userData'), 'fast-storage.log');
  fs.appendFileSync(logFile, `${new Date().toISOString()} - ${message}\n`);

  return true;
});

// Handle App termination
ipcMain.on('quit', (event, code) => {
  app.exit(code);
});
