import { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { api } from '../../lib/api';

export default function NotesIndex() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [hasNotes, setHasNotes] = useState(false);

  useEffect(() => {
    async function checkNotes() {
      try {
        const data = await api('/api/notes');
        
        if (data && Array.isArray(data) && data.length > 0) {
          // Redirect to first note
          router.replace(`/notes/${data[0].id}`);
        } else {
          // No notes, show empty state
          setHasNotes(false);
          setLoading(false);
        }
      } catch (err) {
        console.error('Failed to load notes:', err);
        if (err.message.toLowerCase().includes('unauthorized')) {
          router.push('/');
        } else {
          setLoading(false);
        }
      }
    }

    checkNotes();
  }, []);

  const handleCreateNote = async () => {
    try {
      const newNote = await api('/api/notes', {
        method: 'POST',
        body: JSON.stringify({
          title: 'Untitled Note',
          content: '',
          shared: false,
          favorite: false
        }),
      });
      router.push(`/notes/${newNote.id}`);
    } catch (err) {
      console.error(err);
      alert('Gagal membuat note');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: '#0E2148' }}>
        <div className="text-center">
          <svg className="w-20 h-20 mx-auto mb-4 text-gray-600 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p className="text-gray-400 text-xl font-medium">Loading notes...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center" style={{ background: '#0E2148' }}>
      <div className="text-center">
        <svg className="w-32 h-32 mx-auto mb-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h1 className="text-white text-3xl font-bold mb-3">No Notes Yet</h1>
        <p className="text-gray-400 text-lg mb-8">Create your first note to get started</p>
        <button
          onClick={handleCreateNote}
          className="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white text-base font-medium rounded-lg transition-all shadow-lg"
        >
          + Create First Note
        </button>
      </div>
    </div>
  );
}