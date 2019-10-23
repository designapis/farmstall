const express = require('express')
const app = express()
const path = require('path')

// Middleware
const bodyParser = require('body-parser')
const morgan = require('morgan')

// Configs
const PORT = process.env.NODE_PORT || 3001

// Basic middleware
app.use(morgan('common'))
app.use(bodyParser.json())

// Static assets
app.use(express.static(path.join(__dirname, 'build')));

// Routes
app.get('/health', (req,res,next) => {
  res.send('healthy')
})

app.post('/debug', (req,res,next) => {
  res.send({
    body: req.body,
    headers: req.headers,
    query: req.query,
  })
})

// SPA
app.get((req,res) => {
  res.sendFile(path.join(__dirname, 'build', 'index.html'))
})

// Start
const server = app.listen(PORT, () => {
  console.log('Server listening on http://' + server.address().address + ':' + server.address().port)
})
