import { useState, type FormEvent } from 'react'
import { login } from '../services/api'
import { errorMessage } from '../utils'

interface LoginScreenProps { onLogin: (session: Awaited<ReturnType<typeof login>>) => Promise<void> }

export function LoginScreen({ onLogin }: LoginScreenProps) {
  const [email, setEmail] = useState('')
  const [senha, setSenha] = useState('')
  const [feedback, setFeedback] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setFeedback(''); setLoading(true)
    try { await onLogin(await login(email, senha)) } catch (error) { setFeedback(errorMessage(error, 'Não foi possível entrar.')) } finally { setLoading(false) }
  }

  return <section className="login-view">
    <div className="login-intro"><span className="kicker">Acesso restrito</span><h2>Organize experiências que ficam.</h2><p>Entre no seu painel para criar, acompanhar e manter seus eventos em um só lugar.</p><div className="intro-line"><span /><span /><span /></div></div>
    <form className="form-panel" onSubmit={submit}>
      <div className="form-heading"><span className="form-icon">↳</span><div><h3>Entrar no painel</h3><p>Use as credenciais da sua conta.</p></div></div>
      <label htmlFor="login-email">E-mail</label><input id="login-email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} placeholder="empresa@exemplo.com" required />
      <label htmlFor="login-password">Senha</label><input id="login-password" type="password" value={senha} onChange={(event) => setSenha(event.target.value)} placeholder="Sua senha" required />
      <button className="primary-button full-width" disabled={loading}>{loading ? 'Entrando...' : 'Entrar'} <b>→</b></button>
      {feedback && <p className="feedback">{feedback}</p>}
    </form>
  </section>
}
