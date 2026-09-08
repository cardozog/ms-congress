import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { EventForm } from '../components/EventForm'
import { createEvent, getEvent, getProfile, getTicketTypes, updateEvent } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventDetails, EventPayload, TicketType } from '../types'
import { errorMessage } from '../utils'

export function EventFormPage() {
  const { id } = useParams(); const { session } = useSession(); const navigate = useNavigate(); const [event, setEvent] = useState<EventDetails | null>(null); const [organizerId, setOrganizerId] = useState(0); const [ticketTypes, setTicketTypes] = useState<TicketType[]>([]); const [feedback, setFeedback] = useState(''); const [loading, setLoading] = useState(false)
  useEffect(() => { if (!session) return; const currentSession = session; async function load() { try { const [profile, types] = await Promise.all([getProfile(currentSession), getTicketTypes(currentSession)]); if (!profile.pessoaJuridica?.id) throw new Error('A conta não possui um cadastro de pessoa jurídica associado.'); setOrganizerId(profile.pessoaJuridica.id); setTicketTypes(types); if (id) setEvent(await getEvent(currentSession, Number(id))) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível carregar o formulário.')) } } void load() }, [session, id])
  async function save(payload: EventPayload) { setLoading(true); setFeedback(''); try { if (event) await updateEvent(session!, event.id, payload); else await createEvent(session!, payload); navigate('/eventos', { replace: true }) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível salvar o evento.')) } finally { setLoading(false) } }
  return <EventForm event={event} ticketTypes={ticketTypes} organizerId={organizerId} loading={loading} feedback={feedback} onSubmit={save} onCancel={() => navigate('/eventos')} />
}
