import { useEffect, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { getOrganizerEventReport, getProfile } from '../services/api'
import { useSession } from '../session/useSession'
import type { EventReport } from '../types'
import { errorMessage, formatDate } from '../utils'

const currency = new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' })

export function EventReportPage() {
  const { session } = useSession()
  const [organizerId, setOrganizerId] = useState<number | null>(null)
  const [events, setEvents] = useState<EventReport[]>([])
  const [feedback, setFeedback] = useState('')
  const [loadingOrganizer, setLoadingOrganizer] = useState(true)
  const [loading, setLoading] = useState(false)
  const [hasSearched, setHasSearched] = useState(false)
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')

  useEffect(() => {
    if (!session) return
    let active = true
    async function loadOrganizer() {
      try {
        const profile = await getProfile(session)
        const id = profile.pessoaJuridica?.id
        if (!id) throw new Error('O perfil não possui um ID de organizador.')
        if (active) setOrganizerId(id)
      } catch (error) {
        if (active) setFeedback(errorMessage(error, 'Não foi possível carregar os dados do organizador.'))
      } finally {
        if (active) setLoadingOrganizer(false)
      }
    }
    void loadOrganizer()
    return () => { active = false }
  }, [session])

  const totals = events.reduce((summary, event) => ({
    tickets: summary.tickets + event.ingressosVendidos,
    attendance: summary.attendance + event.presencas,
    revenue: summary.revenue + event.valorArrecadado,
  }), { tickets: 0, attendance: 0, revenue: 0 })

  async function searchReport(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!session || !organizerId) return
    if (dateFrom && dateTo && dateFrom > dateTo) {
      setFeedback('A data final deve ser igual ou posterior à data inicial.')
      return
    }

    setLoading(true)
    setFeedback('')
    setHasSearched(false)
    setEvents([])
    try {
      setEvents(await getOrganizerEventReport(session, organizerId, dateFrom, dateTo))
      setHasSearched(true)
    } catch (error) {
      setFeedback(errorMessage(error, 'Não foi possível carregar o relatório.'))
    } finally {
      setLoading(false)
    }
  }

  return <section className="view-content">
    <div className="section-toolbar">
      <div>
        <span className="kicker">Desempenho</span>
        <h2>Relatório de eventos</h2>
        <p className="section-description">Acompanhe ingressos vendidos, presenças e arrecadação dos seus eventos.</p>
      </div>
    </div>
    {loadingOrganizer && <p className="section-description" role="status">Carregando dados do organizador...</p>}
    {feedback && <p className="feedback" role="alert">{feedback}</p>}
    <form className="report-filters" onSubmit={searchReport}>
      <label>
        Data de
        <input type="date" value={dateFrom} max={dateTo || undefined} onChange={event => setDateFrom(event.target.value)} disabled={loading || loadingOrganizer} />
      </label>
      <label>
        Data até
        <input type="date" value={dateTo} min={dateFrom || undefined} onChange={event => setDateTo(event.target.value)} disabled={loading || loadingOrganizer} />
      </label>
      <button className="primary-button" type="submit" disabled={loading || loadingOrganizer || !organizerId}>
        {loading ? 'Pesquisando...' : 'Pesquisar'}
      </button>
    </form>
    {loading ? <p className="section-description" role="status">Consultando relatório...</p> : hasSearched ? <>
      <div className="report-summary">
        <article className="report-metric"><span>Eventos</span><strong>{events.length}</strong></article>
        <article className="report-metric"><span>Ingressos vendidos</span><strong>{totals.tickets}</strong></article>
        <article className="report-metric"><span>Presenças confirmadas</span><strong>{totals.attendance}</strong></article>
        <article className="report-metric"><span>Valor arrecadado</span><strong>{currency.format(totals.revenue)}</strong></article>
      </div>
      {!events.length ? <div className="empty-state report-empty"><h3>Nenhum evento encontrado</h3><p>Não há eventos para o período informado.</p></div> :
        <div className="report-table-container" role="region" aria-label="Relatório de eventos" tabIndex={0}>
          <table className="report-table">
            <thead>
              <tr>
                <th scope="col">Evento</th>
                <th scope="col">Período</th>
                <th scope="col" className="report-number">Ingressos vendidos</th>
                <th scope="col" className="report-number">Presenças</th>
                <th scope="col" className="report-number">Valor arrecadado</th>
                <th scope="col">Ações</th>
              </tr>
            </thead>
            <tbody>
              {events.map(event => <tr key={event.id}>
                <th scope="row" className="report-event-name">{event.nome}</th>
                <td>{formatDate(event.dataInicio)} – {formatDate(event.dataFim)}</td>
                <td className="report-number">{event.ingressosVendidos}</td>
                <td className="report-number">{event.presencas}</td>
                <td className="report-number">{currency.format(event.valorArrecadado)}</td>
                <td><Link className="secondary-button" to={`/eventos/${event.id}`}>Ver evento</Link></td>
              </tr>)}
            </tbody>
          </table>
        </div>}
    </> : null}
  </section>
}
