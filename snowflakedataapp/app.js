import express from 'express';
import db from './db.js';

const app = express();

app.get('/api/:table', async (req, res) => {
  const table = req.params.table;
  if (!db.tables.includes(table)) throw new Error(`Table ${table} not in list of definitons`);
  const result = await db.execute(`SELECT * FROM "${table}" LIMIT 1000`);
  res.send(result);
});

export default app;