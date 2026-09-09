import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { EventCard } from '../components/EventCard'
import { ConfirmModal } from '../components/ConfirmModal'
import { closeEvent, getEvents, getProfile, publishEvent } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventSummary, Profile } from '../types'
import { errorMessage } from '../utils'

type PendingAction = { kind: 'publish' | 'close'; id: number }

export function EventsPage() {
  const { session } = useSession()
  const navigate = useNavigate()
  const [events, setEvents] = useState<EventSummary[]>([])
  const [profile, setProfile] = useState<Profile | null>(null)
  const [feedback, setFeedback] = useState('')
  const [publishingId, setPublishingId] = useState<number | null>(null)
  const [closingId, setClosingId] = useState<number | null>(null)
  const [pendingAction, setPendingAction] = useState<PendingAction | null>(null)
  useEffect(() => { if (!session) return; const currentSession = session; async function load() { try { const nextProfile = await getProfile(currentSession); setProfile(nextProfile); const id = nextProfile.pessoaJuridica?.id; if (!id) throw new Error('A conta não possui um cadastro de pessoa jurídica associado.'); setEvents(await getEvents(currentSession, id)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível carregar os eventos.')) } } void load() }, [session])
  async function publish(id: number) { if (!session) return; setPublishingId(id); setFeedback(''); try { const published = await publishEvent(session, id); setEvents((current) => current.map((event) => event.id === id ? { ...event, publicado: published.publicado } : event)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível publicar o evento.')) } finally { setPublishingId(null); setPendingAction(null) } }
  async function close(id: number) { if (!session) return; setClosingId(id); setFeedback(''); try { const closed = await closeEvent(session, id); setEvents((current) => current.map((event) => event.id === id ? { ...event, dataFim: closed.dataFim, encerrado: closed.encerrado } : event)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível encerrar o evento.')) } finally { setClosingId(null); setPendingAction(null) } }
  function confirmAction() { if (!pendingAction) return; if (pendingAction.kind === 'publish') void publish(pendingAction.id); else void close(pendingAction.id) }
  const actionIsLoading = pendingAction?.kind === 'publish' ? publishingId === pendingAction.id : closingId === pendingAction?.id
  return <section className="view-content"><div className="section-toolbar"><div><span className="kicker">Visão geral</span><p className="section-description">{events.length} {events.length === 1 ? 'evento encontrado' : 'eventos encontrados'} para sua organização.</p></div><Link className="primary-button" to="/eventos/novo">+ Novo evento</Link></div>{feedback && <p className="feedback">{feedback}</p>}<div className="events-list">{events.length ? events.map((event) => <EventCard key={event.id} event={event} onSelect={(eventId) => navigate(`/eventos/${eventId}`)} onPublish={(eventId) => setPendingAction({ kind: 'publish', id: eventId })} onClose={(eventId) => setPendingAction({ kind: 'close', id: eventId })} publishing={publishingId === event.id} closing={closingId === event.id} />) : <div className="empty-state"><h3>Seu calendário começa aqui.</h3><p>Crie o primeiro evento da sua organização para começar.</p></div>}</div>{profile && <span className="page-data" data-organizer-id={profile.pessoaJuridica?.id} />}{pendingAction && <ConfirmModal title={pendingAction.kind === 'publish' ? 'Publicar evento?' : 'Encerrar evento?'} message={pendingAction.kind === 'publish' ? 'O evento ficará disponível como publicado e não poderá mais ser editado.' : 'O evento será encerrado imediatamente e não poderá voltar ao estado ativo.'} confirmLabel={pendingAction.kind === 'publish' ? 'Publicar' : 'Encerrar'} tone={pendingAction.kind} loading={actionIsLoading} onConfirm={confirmAction} onCancel={() => setPendingAction(null)} />}</section>
}
