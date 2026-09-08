import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { EventCard } from '../components/EventCard'
import { getEvents, getProfile } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventSummary, Profile } from '../types'
import { errorMessage } from '../utils'

export function EventsPage() {
  const { session } = useSession()
  const navigate = useNavigate()
  const [events, setEvents] = useState<EventSummary[]>([])
  const [profile, setProfile] = useState<Profile | null>(null)
  const [feedback, setFeedback] = useState('')
  useEffect(() => { if (!session) return; const currentSession = session; async function load() { try { const nextProfile = await getProfile(currentSession); setProfile(nextProfile); const id = nextProfile.pessoaJuridica?.id; if (!id) throw new Error('A conta não possui um cadastro de pessoa jurídica associado.'); setEvents(await getEvents(currentSession, id)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível carregar os eventos.')) } } void load() }, [session])
  return <section className="view-content"><div className="section-toolbar"><div><span className="kicker">Visão geral</span><p className="section-description">{events.length} {events.length === 1 ? 'evento encontrado' : 'eventos encontrados'} para sua organização.</p></div><Link className="primary-button" to="/eventos/novo">+ Novo evento</Link></div>{feedback && <p className="feedback">{feedback}</p>}<div className="events-list">{events.length ? events.map((event) => <EventCard key={event.id} event={event} onSelect={(eventId) => navigate(`/eventos/${eventId}`)} />) : <div className="empty-state"><h3>Seu calendário começa aqui.</h3><p>Crie o primeiro evento da sua organização para começar.</p></div>}</div>{profile && <span className="page-data" data-organizer-id={profile.pessoaJuridica?.id} />}</section>
}
