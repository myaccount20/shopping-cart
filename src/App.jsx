import { useState, useEffect } from 'react'
import Login from './components/Login'
import ItemsList from './components/ItemsList'

const API_URL = 'http://localhost:8080'

function App() {
  const [token, setToken] = useState(null)

  useEffect(() => {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      setToken(savedToken)
    }
  }, [])

  const handleLogin = (newToken) => {
    setToken(newToken)
    localStorage.setItem('token', newToken)
  }

  const handleLogout = () => {
    setToken(null)
    localStorage.removeItem('token')
  }

  if (!token) {
    return <Login onLogin={handleLogin} apiUrl={API_URL} />
  }

  return <ItemsList token={token} apiUrl={API_URL} onLogout={handleLogout} />
}

export default App
