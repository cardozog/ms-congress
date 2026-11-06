import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { QRCodeSVG } from 'qrcode.react'
import { getEvent, getEventCheckins } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventCheckin, EventDetails } from '../types'
import { errorMessage, formatDate } from '../utils'

export function EventCheckinsPage() {
  const { id } = useParams()
  const { session } = useSession()
  const [event, setEvent] = useState<EventDetails | null>(null)
  const [records, setRecords] = useState<EventCheckin[]>([])
  const [search, setSearch] = useState('')
  const [feedback, setFeedback] = useState('')
  const [selected, setSelected] = useState<EventCheckin | null>(null)

  useEffect(() => {
    if (!session || !id) return
    const eventId = Number(id)
    void Promise.all([getEvent(session, eventId), getEventCheckins(session, eventId)])
      .then(([loadedEvent, response]) => { setEvent(loadedEvent); setRecords(response.checkins) })
      .catch(error => setFeedback(errorMessage(error, 'Não foi possível carregar os check-ins.')))
  }, [session, id])

  async function searchCheckins() {
    if (!session || !id) return
    try {
      setFeedback('')
      setRecords((await getEventCheckins(session, Number(id), search)).checkins)
    } catch (error) {
      setFeedback(errorMessage(error, 'Não foi possível buscar os check-ins.'))
    }
  }

  return <section className="checkins-view">
    <Link className="back-button no-print" to={`/eventos/${id}`}>← Voltar para o evento</Link>
    <div className="section-toolbar no-print"><div><p className="eyebrow">Presença confirmada</p><h2>{event?.nome ?? 'Check-ins do evento'}</h2><p className="section-description">{records.length} participante{records.length === 1 ? '' : 's'} fizeram check-in.</p></div></div>
    <div className="checkin-search no-print"><input value={search} onChange={event => setSearch(event.target.value)} onKeyDown={event => { if (event.key === 'Enter') void searchCheckins() }} placeholder="Buscar por nome, CPF ou e-mail" aria-label="Buscar participante" /><button className="primary-button" type="button" onClick={() => void searchCheckins()}>Buscar</button></div>
    {feedback && <div className="feedback no-print" role="alert">{feedback}</div>}
    {!records.length ? <div className="empty-state no-print"><h3>Nenhum check-in encontrado</h3><p>Os participantes aparecerão aqui após confirmarem presença.</p></div> : <div className="checkin-list no-print">{records.map(record => <article className="checkin-card" key={record.ingressoId}><div><span className="kicker">Check-in confirmado</span><h3>{record.nome}</h3><p>{record.tipoIngresso} · CPF {record.cpf}</p><small>{record.email}{record.checkinEm ? ` · ${formatDate(record.checkinEm, true)}` : ''}</small></div><button className="secondary-button" type="button" onClick={() => setSelected(record)}>Imprimir credencial</button></article>)}</div>}
    {selected && <section className="credential-print" aria-label="Credencial do participante"><div className="credential-print-top">{event?.eventoLogo && <img src={event.eventoLogo} alt="" />}<span>MS Congress</span></div><div className="credential-print-body"><p className="kicker">Credencial de participante</p><h1>{event?.nome ?? selected.evento}</h1><h2>{selected.nome}</h2><p className="credential-role">{selected.tipoIngresso}</p><div className="credential-print-details"><span>Documento<strong>{selected.cpf}</strong></span><span>E-mail<strong>{selected.email}</strong></span><span>Nascimento<strong>{selected.dataNascimento ?? 'Não informado'}</strong></span><span>Check-in<strong>{selected.checkinEm ? formatDate(selected.checkinEm, true) : 'Confirmado'}</strong></span></div><QRCodeSVG value={`ms-congress:credential:${selected.codigo}`} size={150} level="M" includeMargin /><code>{selected.codigo}</code></div><div className="credential-actions no-print"><button className="secondary-button" type="button" onClick={() => setSelected(null)}>Fechar</button><button className="primary-button" type="button" onClick={() => window.print()}>Imprimir credencial</button></div></section>}
  </section>
}
