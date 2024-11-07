import snowflake from 'snowflake-sdk';
import config from './config.js';
import data from './data/northwind.min.json' with {type: 'json'};

console.log({ config });

class Db {
  connection = null;
  definitions = {};
  tables = [];
  
  connect = async () => new Promise((resolve, reject) => {
    snowflake.createConnection({
      account: config.SNOWFLAKE_ACCOUNT,
      username: config.SNOWFLAKE_USER,
      password: config.SNOWFLAKE_PASSWORD,
      application: config.SNOWFLAKE_APP,
      authenticator: 'SNOWFLAKE',
      role: 'DEFAULT_WH',
      database: config.SNOWFLAKE_DATABASE,
      schema: config.SNOWFLAKE_SCHEMA,
    }).connect((err, conn) => {
      if (err) {
        console.error('Unable to connect: ' + err.message);
        reject(err.message);
      }
      else {
        console.log('Successfully connected to Snowflake....', conn.isUp());
        this.connection = conn;
        resolve(this);
      }
    });
  });

  execute = (sqlText, opts) => new Promise((resolve, reject) => {
    console.log(`SQL> ${sqlText}`);
    if (!this.connection) throw new Error('SQL Connection not ready..');
    this.connection.execute({
      ...opts,
      sqlText,
      complete: (err, stmt, rows) => {
        if (err) {
          console.error(err);
          reject(err);
        } else {
          resolve(rows);
        }
      }
    });
  });

  preload = async (schema) => {
    console.log('preloading schema', schema, Object.keys(data));
    await this.execute(`CREATE SCHEMA IF NOT EXISTS ${schema}`);
    await this.execute(`USE SCHEMA ${schema}`);
    for(const table of Object.keys(data)) {
      const record = {...data[table][0]};
      console.log({record});
      const definition = this.definitions[table] = [];
      for (const col of Object.keys(record)) {
        const val = record[col];
        let type = 'VARCHAR';
        if (typeof val === 'number' || String(col).toLowerCase().endsWith('id') || String(col).toLowerCase().includes('unit')) {
          type = 'INT';
        }
        definition.push(`"${col}" ${type}`);
      }
      await this.execute(`CREATE TABLE IF NOT EXISTS "${table}" (${definition.join(', ')})`);
      await this.execute(`INSERT INTO "${table}" VALUES (${new Array(definition.length).fill('?')})`, {
        binds: data[table].map(obj => Object.values(obj)),
      });
      this.tables.push(table);
    }
  }
};

export default new Db();
