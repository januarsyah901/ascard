import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { 
  LayoutDashboard, 
  FileText, 
  Plus, 
  LogOut, 
  Trash2, 
  Edit2, 
  CheckCircle2, 
  XCircle,
  ExternalLink,
  ChevronRight
} from 'lucide-react';

const AdminDashboard = () => {
  const { user, token, logout } = useAuth();
  const navigate = useNavigate();
  const [articles, setArticles] = useState([]);
  const [isEditing, setIsEditing] = useState(false);
  const [currentArticle, setCurrentArticle] = useState({ title: '', content: '', category: '', published: false });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!token) navigate('/admin/login');
    fetchArticles();
  }, [token]);

  const fetchArticles = async () => {
    try {
      const res = await fetch('http://localhost:5001/api/articles');
      const data = await res.json();
      setArticles(data);
    } catch (err) {
      console.error('Error fetching articles:', err);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const method = currentArticle.id ? 'PUT' : 'POST';
    const url = currentArticle.id 
      ? `http://localhost:5001/api/articles/${currentArticle.id}` 
      : 'http://localhost:5001/api/articles';

    try {
      const res = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(currentArticle)
      });

      if (res.ok) {
        setIsEditing(false);
        setCurrentArticle({ title: '', content: '', category: '', published: false });
        fetchArticles();
      }
    } catch (err) {
      console.error('Error saving article:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this article?')) return;
    try {
      const res = await fetch(`http://localhost:5001/api/articles/${id}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) fetchArticles();
    } catch (err) {
      console.error('Error deleting article:', err);
    }
  };

  const startEdit = (article) => {
    setCurrentArticle(article);
    setIsEditing(true);
  };

  return (
    <div className="min-h-screen bg-background text-on-surface flex font-body">
      {/* Sidebar */}
      <aside className="w-72 bg-surface-container-lowest border-r border-outline flex flex-col">
        <div className="p-8 border-b border-outline-variant/30 mb-8">
          <div className="bg-primary text-on-primary-fixed px-2 py-1 inline-block font-label text-[9px] font-black tracking-[0.2em] uppercase mb-4">
            System Control
          </div>
          <h1 className="font-headline font-black text-2xl tracking-tighter uppercase leading-none">
            Clinical <br /> Dashboard
          </h1>
        </div>

        <nav className="flex-1 px-4 space-y-1">
          <button className="w-full flex items-center justify-between bg-primary text-on-primary-fixed px-6 py-4 rounded-none transition-all group">
            <div className="flex items-center gap-4">
              <FileText className="w-4 h-4" />
              <span className="font-label font-bold text-[11px] uppercase tracking-widest">Articles</span>
            </div>
            <ChevronRight className="w-4 h-4 opacity-0 group-hover:opacity-100 transition-all" />
          </button>
          
          <a 
            href="/knowledge-base" 
            target="_blank" 
            className="w-full flex items-center justify-between hover:bg-surface-container-high px-6 py-4 rounded-none transition-all group mt-4 border border-outline-variant/30"
          >
            <div className="flex items-center gap-4">
              <ExternalLink className="w-4 h-4" />
              <span className="font-label font-bold text-[11px] uppercase tracking-widest">View Public Site</span>
            </div>
          </a>
        </nav>

        <div className="p-8 mt-auto border-t border-outline-variant/30">
          <button 
            onClick={logout}
            className="flex items-center gap-3 text-outline hover:text-primary transition-colors font-label font-bold text-[10px] uppercase tracking-widest"
          >
            <LogOut className="w-4 h-4" />
            <span>Terminate Session</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 flex flex-col h-screen overflow-hidden">
        <header className="px-12 py-8 border-b border-outline-variant/30 flex justify-between items-end bg-surface">
          <div>
            <div className="flex items-center gap-2 mb-2">
              <span className="w-2 h-2 bg-primary"></span>
              <span className="font-label font-bold text-[10px] uppercase tracking-[0.3em] text-outline">Content Repository</span>
            </div>
            <h2 className="font-headline font-black text-4xl tracking-tighter uppercase">Knowledge Base</h2>
          </div>
          {!isEditing && (
            <button 
              onClick={() => { setIsEditing(true); setCurrentArticle({ title: '', content: '', category: '', published: false }); }}
              className="flex items-center gap-3 bg-primary hover:bg-primary-fixed-dim text-on-primary-fixed px-8 py-4 rounded-none font-label font-black text-xs uppercase tracking-[0.2em] transition-all"
            >
              <Plus className="w-4 h-4" />
              <span>Create Insight</span>
            </button>
          )}
        </header>

        <div className="flex-1 p-12 overflow-y-auto">
          {isEditing ? (
            <div className="max-w-4xl animate-in slide-in-from-bottom-4 duration-500">
              <div className="flex items-center gap-4 mb-8">
                <button onClick={() => setIsEditing(false)} className="text-outline hover:text-primary transition-colors">
                  <span className="font-label font-bold text-[10px] uppercase tracking-widest border-b border-outline">Back to Inventory</span>
                </button>
                <span className="text-outline-variant">/</span>
                <span className="font-label font-bold text-[10px] uppercase tracking-widest text-primary">
                  {currentArticle.id ? 'Modify Existing' : 'Initialize New Entry'}
                </span>
              </div>

              <form onSubmit={handleSubmit} className="space-y-10">
                <div className="grid grid-cols-2 gap-10">
                  <div className="space-y-3">
                    <label className="block font-label font-bold text-[11px] uppercase tracking-widest text-on-surface-variant">Entry Title</label>
                    <input 
                      type="text" 
                      value={currentArticle.title}
                      onChange={(e) => setCurrentArticle({...currentArticle, title: e.target.value})}
                      className="w-full bg-surface-container-low border-b-2 border-outline-variant p-4 rounded-none focus:border-primary outline-none font-headline font-bold text-xl tracking-tight transition-all"
                      placeholder="ENTER TITLE..."
                      required
                    />
                  </div>
                  <div className="space-y-3">
                    <label className="block font-label font-bold text-[11px] uppercase tracking-widest text-on-surface-variant">Category Tag</label>
                    <input 
                      type="text" 
                      value={currentArticle.category}
                      onChange={(e) => setCurrentArticle({...currentArticle, category: e.target.value})}
                      className="w-full bg-surface-container-low border-b-2 border-outline-variant p-4 rounded-none focus:border-primary outline-none font-label font-bold text-xs uppercase tracking-widest transition-all"
                      placeholder="PSYCHOLOGY / MATH / FINANCE"
                    />
                  </div>
                </div>

                <div className="space-y-3">
                  <label className="block font-label font-bold text-[11px] uppercase tracking-widest text-on-surface-variant">Clinical Data (Markdown)</label>
                  <textarea 
                    value={currentArticle.content}
                    onChange={(e) => setCurrentArticle({...currentArticle, content: e.target.value})}
                    className="w-full bg-surface-container-lowest border border-outline-variant p-6 rounded-none focus:border-primary outline-none min-h-[400px] font-body text-base leading-relaxed transition-all"
                    placeholder="Enter analytical content..."
                    required
                  />
                </div>

                <div className="flex items-center gap-4 bg-surface-container-low p-6 border border-outline-variant/30">
                  <input 
                    type="checkbox" 
                    id="published"
                    checked={currentArticle.published}
                    onChange={(e) => setCurrentArticle({...currentArticle, published: e.target.checked})}
                    className="w-6 h-6 rounded-none border-primary bg-transparent text-primary focus:ring-0 appearance-none border-2 checked:bg-primary relative after:content-[''] after:hidden checked:after:block after:absolute after:left-1.5 after:top-0.5 after:w-1.5 after:h-3 after:border-white after:border-r-2 after:border-b-2 after:rotate-45"
                  />
                  <label htmlFor="published" className="font-label font-bold text-[10px] uppercase tracking-[0.2em] text-on-surface cursor-pointer">Authorize for Public Access</label>
                </div>

                <div className="flex gap-6 pt-4">
                  <button 
                    type="submit" 
                    disabled={loading}
                    className="bg-primary hover:bg-primary-fixed-dim text-on-primary-fixed px-12 py-5 rounded-none font-label font-black text-xs uppercase tracking-[0.3em] transition-all disabled:opacity-50"
                  >
                    {loading ? 'COMMITTING...' : 'SAVE TO REPOSITORY'}
                  </button>
                  <button 
                    type="button" 
                    onClick={() => setIsEditing(false)}
                    className="border border-outline hover:bg-surface-container-high px-12 py-5 rounded-none font-label font-black text-xs uppercase tracking-[0.3em] transition-all"
                  >
                    ABORT
                  </button>
                </div>
              </form>
            </div>
          ) : (
            <div className="grid gap-px bg-outline-variant/30 border border-outline-variant/30">
              {articles.length === 0 ? (
                <div className="text-center py-32 bg-surface">
                  <div className="w-16 h-16 border-2 border-dashed border-outline-variant mx-auto mb-6 flex items-center justify-center">
                    <FileText className="w-6 h-6 text-outline-variant" />
                  </div>
                  <p className="font-label font-bold text-[10px] uppercase tracking-[0.2em] text-outline">Inventory is empty. Initializing required.</p>
                </div>
              ) : (
                articles.map(article => (
                  <div key={article.id} className="bg-surface p-8 flex items-center justify-between group hover:bg-surface-container-lowest transition-all">
                    <div className="flex items-center gap-8">
                      <div className="font-label font-black text-xs text-outline-variant/30 w-8">
                        {String(articles.indexOf(article) + 1).padStart(2, '0')}
                      </div>
                      <div className={`w-3 h-3 ${article.published ? 'bg-primary' : 'border border-outline-variant'}`}></div>
                      <div>
                        <h4 className="font-headline font-black text-xl tracking-tighter uppercase mb-2 group-hover:text-primary transition-colors">{article.title}</h4>
                        <div className="flex items-center gap-4">
                          <span className="bg-surface-container-highest px-3 py-1 font-label font-bold text-[9px] uppercase tracking-widest text-on-surface-variant">
                            {article.category || 'GENERAL'}
                          </span>
                          <span className="font-label text-[9px] uppercase tracking-[0.2em] text-outline">
                            Updated: {new Date(article.updatedAt).toLocaleDateString()}
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-all translate-x-4 group-hover:translate-x-0">
                      <button 
                        onClick={() => startEdit(article)}
                        className="p-4 hover:bg-primary hover:text-on-primary-fixed transition-all"
                        title="Modify"
                      >
                        <Edit2 className="w-4 h-4" />
                      </button>
                      <button 
                        onClick={() => handleDelete(article.id)}
                        className="p-4 hover:bg-error hover:text-on-error transition-all"
                        title="Erase"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                      <button className="p-4 hover:bg-surface-container-highest transition-all" title="Inspect">
                        <ExternalLink className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};

export default AdminDashboard;
