import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getMyEvent } from '../services/api'
import { useSession } from '../session/useSession'
import type { MyEvent } from '../types'
import { errorMessage } from '../utils'
import { QRCodeSVG } from 'qrcode.react'

export function CongressistaMyEventDetailPage() {
  const { session } = useSession()
  const { id } = useParams()
  const [event, setEvent] = useState<MyEvent | null>(null)
  const [error, setError] = useState('')
  const [visibleQRCode, setVisibleQRCode] = useState<number | null>(null)
  useEffect(() => { if (session && id) void getMyEvent(session, Number(id)).then(setEvent).catch(error => setError(errorMessage(error, 'Não foi possível carregar os ingressos.'))) }, [session, id])
  useEffect(() => {
    if (visibleQRCode === null) return
    const closeOnEscape = (event: KeyboardEvent) => { if (event.key === 'Escape') setVisibleQRCode(null) }
    document.addEventListener('keydown', closeOnEscape)
    return () => document.removeEventListener('keydown', closeOnEscape)
  }, [visibleQRCode])
  function downloadQRCode(ticketId: number, code: string) {
    const svg = document.getElementById(`qr-code-${ticketId}`)
    if (!svg) return
    const blob = new Blob([svg.outerHTML], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `credencial-${code}.svg`
    link.click()
    URL.revokeObjectURL(url)
  }

  const selectedTicket = event?.ingressos.find(ticket => ticket.id === visibleQRCode)

  return <div className="app-shell congressista-shell"><main className="main-content"><section className="my-events-view">  <Link className="secondary-button" to="/meus-eventos">← Meus Eventos</Link>{error ? <div className="empty-state"><p>{error}</p></div> : !event ? <div className="empty-state"><p>Carregando ingressos...</p></div> : <><span className="kicker">Credenciais digitais</span><h1>{event.nome}</h1><p>{event.quantidade} ingressos adquiridos</p><div className="credential-list">{event.ingressos.map(ticket => <div className="credential-card" key={ticket.id}><strong>{ticket.ingresso}</strong><span>Compra confirmada</span><code>{ticket.codigo}</code><button className="secondary-button" type="button" onClick={() => setVisibleQRCode(ticket.id)}>Visualizar QR Code</button></div>)}</div></>}</section></main>{selectedTicket && <div className="qr-modal-backdrop" role="presentation" onClick={() => setVisibleQRCode(null)}><section className="qr-modal" role="dialog" aria-modal="true" aria-labelledby="qr-modal-title" onClick={event => event.stopPropagation()}><button className="qr-modal-close" type="button" aria-label="Fechar QR Code" onClick={() => setVisibleQRCode(null)}>×</button><span className="kicker">Credencial digital</span><h2 id="qr-modal-title">{selectedTicket.ingresso}</h2><QRCodeSVG id={`qr-code-${selectedTicket.id}`} value={selectedTicket.qrCode ?? selectedTicket.codigo} size={220} level="M" includeMargin aria-label={`QR Code do ingresso ${selectedTicket.codigo}`} /><code>{selectedTicket.codigo}</code><small>Apresente este QR Code na validação do evento.</small><button className="primary-button" type="button" onClick={() => downloadQRCode(selectedTicket.id, selectedTicket.codigo)}>Baixar imagem</button></section></div>}</div>
}
