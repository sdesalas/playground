import db from './db.js';
import config from './config.js';
import app from './app.js';

(async () => {
  console.log('starting..');
  await db.connect();
  console.log('Using Northwind...');
  await db.execute('USE DATABASE Northwind');
  const schema = 'schema' + String(new Date().getTime());
  await db.preload(schema);
  app.listen(config.PORT, () => console.log(`Listening on port ${config.PORT}`));
})()