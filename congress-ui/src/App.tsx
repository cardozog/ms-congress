import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AppLayout } from './components/AppLayout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { EventDetailPage } from './pages/EventDetailPage'
import { EventFormPage } from './pages/EventFormPage'
import { EventsPage } from './pages/EventsPage'
import { LoginPage } from './pages/LoginPage'
import { ProfilePage } from './pages/ProfilePage'
import { SessionProvider } from './session/SessionContext'
import { getStoredSession } from './services/api'
import './App.css'

function App() {
  return <SessionProvider initialSession={getStoredSession()}><BrowserRouter><Routes><Route path="/login" element={<LoginPage />} /><Route element={<ProtectedRoute />}><Route element={<AppLayout />}><Route index element={<Navigate to="/eventos" replace />} /><Route path="eventos" element={<EventsPage />} /><Route path="eventos/novo" element={<EventFormPage />} /><Route path="eventos/:id" element={<EventDetailPage />} /><Route path="eventos/:id/editar" element={<EventFormPage />} /><Route path="perfil" element={<ProfilePage />} /></Route></Route><Route path="*" element={<Navigate to="/eventos" replace />} /></Routes></BrowserRouter></SessionProvider>
}

export default App
