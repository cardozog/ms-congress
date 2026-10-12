import { useState, type FormEvent } from 'react'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { loginCongressista } from '../services/api'
import { useSession } from '../session/useSession'
import { errorMessage } from '../utils'
import { ArrowRight, LogIn } from 'lucide-react'

export function CongressistaLoginPage() {
  const { session, signIn } = useSession()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [senha, setSenha] = useState('')
  const [feedback, setFeedback] = useState('')
  const [loading, setLoading] = useState(false)
  if (session) return <Navigate to={session.papeis.includes('PARTICIPANTE') ? '/congressista/eventos' : '/eventos'} replace />

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setFeedback(''); setLoading(true)
    try { signIn(await loginCongressista(email, senha)); navigate('/congressista/eventos', { replace: true }) }
    catch (error) { setFeedback(errorMessage(error, 'Não foi possível entrar.')) }
    finally { setLoading(false) }
  }

  return <><header className="topbar public-topbar"><a className="brand" href="/congressista/login"><span className="brand-mark">MS</span><span><strong>MS Congress</strong><small>área do congressista</small></span></a></header><section className="login-view"><div className="login-intro"><span className="kicker">Área do congressista</span><h2>Viva o próximo encontro.</h2><p>Entre para descobrir eventos e acompanhar suas inscrições.</p></div><form className="form-panel" onSubmit={submit}><div className="form-heading"><span className="form-icon"><LogIn size={19} /></span><div><h3>Entrar como congressista</h3><p>Use as credenciais da sua conta.</p></div></div><label htmlFor="congressista-email">E-mail</label><input id="congressista-email" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required /><label htmlFor="congressista-password">Senha</label><input id="congressista-password" type="password" value={senha} onChange={(event) => setSenha(event.target.value)} required /><button className="primary-button full-width" disabled={loading}>{loading ? 'Entrando...' : 'Entrar'} <ArrowRight size={16} /></button>{feedback && <p className="feedback">{feedback}</p>}<p className="auth-link">Ainda não tem conta? <Link to="/congressista/cadastro">Cadastre-se</Link></p></form></section></>
}
