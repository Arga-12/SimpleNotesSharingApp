import { api } from '../../lib/api';
import React, { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/router';

export default function NotePage() {
  const router = useRouter();
  const { id } = router.query;
  const [notes, setNotes] = useState([]);
  const [selectedNote, setSelectedNote] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const editorRef = useRef(null);
  const cursorPositionRef = useRef(null);

  // Load all notes on mount
  useEffect(() => {
    async function loadNotes() {
      try {
      const data = await api('/api/notes');

      // Jika data null atau bukan array → anggap data kosong
      if (!data || !Array.isArray(data)) {
        console.warn('Response kosong dari server (null)');
        setNotes([]);
        return;
      }

      if (data.length === 0) {
        console.log('Database tersambung, tapi belum ada note.');
        setNotes([]);
        return;
      }

      // Kalau ada data
      setNotes(data);
      } catch (err) {
        console.error('Failed to load notes:', err);

        // Kalau unauthorized → arahkan ke login
        if (err.message.toLowerCase().includes('unauthorized')) {
          alert('Kamu harus login dulu!');
          router.push('/');
          return;
        }

        // Kalau koneksi gagal / DB down → alert error
        if (
          err.message.toLowerCase().includes('failed to fetch') ||
          err.message.toLowerCase().includes('db error') ||
          err.message.toLowerCase().includes('network')
        ) {
          alert('Gagal terhubung ke server atau database.');
        } else {
          alert('Terjadi kesalahan: ' + err.message);
        }
      }
    }

    loadNotes();
  }, []);

  // Load specific note when ID changes in URL
  useEffect(() => {
    if (!id || !notes.length) return;

    async function loadNoteById() {
      try {
        const noteDetail = await api(`/api/notes/${id}`);
        setSelectedNote(noteDetail);
      } catch (err) {
        console.error('Failed to fetch note detail:', err);
        alert('Gagal memuat note');
        // Redirect to first note if current note not found
        if (notes.length > 0) {
          router.replace(`/notes/${notes[0].id}`);
        } else {
          router.replace('/notes');
        }
      }
    }

    loadNoteById();
  }, [id, notes.length]);

  // Helper function to get note color
  const getNoteColor = (note) => (note.shared ? '#7965C1' : '#E3D095');

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
      setNotes([newNote, ...notes]);
      // Navigate to the new note
      router.push(`/notes/${newNote.id}`);
    } catch (err) {
      console.error(err);
      alert('Gagal membuat note');
    }
  };

  // Navigate to note by ID
  const handleSelectNote = (noteId) => {
    router.push(`/notes/${noteId}`);
  };

  // Save cursor position
  const saveCursorPosition = () => {
    const selection = window.getSelection();
    if (selection.rangeCount > 0) {
      const range = selection.getRangeAt(0);
      const preSelectionRange = range.cloneRange();
      preSelectionRange.selectNodeContents(editorRef.current);
      preSelectionRange.setEnd(range.startContainer, range.startOffset);
      const start = preSelectionRange.toString().length;
      
      cursorPositionRef.current = {
        start: start,
        end: start + range.toString().length
      };
    }
  };

  // Restore cursor position
  const restoreCursorPosition = () => {
    if (!cursorPositionRef.current || !editorRef.current) return;
    
    const { start } = cursorPositionRef.current;
    const selection = window.getSelection();
    const range = document.createRange();
    
    let charCount = 0;
    let node, foundStart = false;
    const nodes = editorRef.current.childNodes;
    for (let i = 0; i < nodes.length; i++) {
      node = nodes[i];
      if (node.nodeType === Node.TEXT_NODE) {
        const nextCharCount = charCount + node.length;
        if (start >= charCount && start <= nextCharCount) {
          range.setStart(node, start - charCount);
          range.setEnd(node, start - charCount);
          foundStart = true;
        }
        charCount = nextCharCount;
      } else {
        if (foundStart) break;
        const childNodes = node.childNodes;
        for (let j = 0; j < childNodes.length; j++) {
          const childNode = childNodes[j];
          if (childNode.nodeType === Node.TEXT_NODE) {
            const nextCharCount = charCount + childNode.length;
            if (start >= charCount && start <= nextCharCount) {
              range.setStart(childNode, start - charCount);
              range.setEnd(childNode, start - charCount);
              foundStart = true;
            }
            charCount = nextCharCount;
          }
        }
      }
    }
    
    if (foundStart) {
      selection.removeAllRanges();
      selection.addRange(range);
    }
  };

  // Update note field
  const updateNoteField = (field, value) => {
    if (!selectedNote) return;
    
    const updatedNote = { ...selectedNote, [field]: value, updatedAt: new Date().toISOString() };
    setSelectedNote(updatedNote);
    setNotes(notes.map(n => n.id === selectedNote.id ? updatedNote : n));
  };

  // Restore cursor after content update
  useEffect(() => {
    if (editorRef.current && cursorPositionRef.current) {
      restoreCursorPosition();
    }
  }, [selectedNote?.content]);

  // Delete note
  const handleDeleteNote = async () => {
    if (!selectedNote) return;
    if (!confirm('Yakin ingin hapus note ini?')) return;

    try {
      await api(`/api/notes/${selectedNote.id}`, { method: 'DELETE' });
      const newNotes = notes.filter(n => n.id !== selectedNote.id);
      setNotes(newNotes);
      // Navigate to first note or back to /notes
      if (newNotes.length > 0) {
        router.push(`/notes/${newNotes[0].id}`);
      } else {
        router.push('/notes');
      }
    } catch (err) {
      console.error(err);
      if (err.message.includes('forbidden')) {
        alert('Kamu tidak bisa menghapus note milik orang lain!');
      } else {
        alert('Gagal menghapus note');
      }
    }
  };


  // Helpers
  const stripHtml = (html) => html ? html.replace(/<[^>]*>/g, '') : '';
  const getWordCount = () => {
    const text = editorRef.current ? (editorRef.current.textContent || '') : (selectedNote?.content ? stripHtml(selectedNote.content) : '');
    return text.trim().split(/\s+/).filter(w => w.length > 0).length;
  };

  // Filter notes
  const filteredNotes = notes.filter(note => {
    const searchLower = searchQuery.toLowerCase();
    return note.title.toLowerCase().includes(searchLower) || 
           note.content.toLowerCase().includes(searchLower);
  });

  // Autosave ke server saat note berubah
  useEffect(() => {
    if (!selectedNote || !selectedNote.id) return;

    const timeout = setTimeout(async () => {
      try {
        await api(`/api/notes/${selectedNote.id}`, {
          method: 'PUT',
          body: JSON.stringify({
            title: selectedNote.title,
            content: selectedNote.content,
            shared: selectedNote.shared,
            favorite: selectedNote.favorite,
          }),
        });
        console.log('Note saved!');
      } catch (err) {
        console.error('Failed to save note:', err);
      }
    }, 800); // delay 0.8 detik biar gak spam

    return () => clearTimeout(timeout);
  }, [selectedNote]);

  return (
    <div className="min-h-screen" style={{ background: '#0E2148' }}>
      {/* Main Content - Google Docs Style Layout */}
      <div className="h-screen flex">
        {/* Left Sidebar - Notes List */}
        <div className="w-80 h-full border-r" style={{ background: 'rgba(14, 33, 72, 0.4)', borderColor: 'rgba(227, 208, 149, 0.15)' }}>
          <div className="h-full flex flex-col">
            {/* Sidebar Header */}
            <div className="p-4 border-b" style={{ borderColor: 'rgba(227, 208, 149, 0.15)' }}>
              <h2 className="text-sm font-bold text-white">All Notes ({notes.length})</h2>
              <button
                onClick={handleCreateNote}
                className="mt-2 w-full py-2 end-auto bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-all"
              >
                + New Note
              </button>
            </div>

            {/* Notes List */}
            <div className="flex-1 overflow-y-auto">
              {filteredNotes.map((note) => (
                <button
                  key={note.id}
                  onClick={() => handleSelectNote(note.id)}
                  className={`w-full text-left p-4 border-b transition-all ${
                    selectedNote?.id === note.id ? 'bg-white/10' : 'hover:bg-white/5'
                  }`}
                  style={{ borderColor: 'rgba(255, 255, 255, 0.05)' }}
                >
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <h3 className="text-sm font-semibold text-white truncate flex-1">
                      {note.title}
                    </h3>
                    <span className="text-[10px] text-gray-400 whitespace-nowrap">
                      @{note.ownerUsername || 'unknown'}
                    </span>
                  </div>
                  <p className="text-xs text-gray-400 line-clamp-2 mb-2">{stripHtml(note.content)}</p>
                  <div className="flex items-center gap-2">
                    <span className="text-[10px] text-gray-500">
                      {new Date(note.updatedAt).toLocaleDateString()}
                    </span>
                  </div>
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Right Editor Area */}
        <div className="flex-1 h-full overflow-hidden flex flex-col">
          {selectedNote ? (
            <>
              {/* Editor Toolbar */}
              <div className="px-8 py-4 border-b flex items-center justify-between" style={{ background: 'rgba(14, 33, 72, 0.4)', borderColor: 'rgba(227, 208, 149, 0.15)' }}>
                <div className="flex items-center gap-3">
                  <button
                    onClick={() => updateNoteField('favorite', !selectedNote.favorite)}
                    className="p-2 rounded-lg hover:bg-white/10 transition-all"
                  >
                    <svg className={`w-5 h-5 ${selectedNote.favorite ? 'text-yellow-400 fill-current' : 'text-gray-400'}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                    </svg>
                  </button>
                </div>
                
                <div className="flex items-center gap-3">
                  <button
                    onClick={handleDeleteNote}
                    className="px-4 py-2 rounded-lg text-sm font-medium text-white transition-all hover:opacity-90"
                    style={{ background: 'rgba(239, 68, 68, 0.5)' }}
                  >
                    Delete Note
                  </button>
                  
                  <button
                    onClick={() => router.push('/profile')}
                    className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 transition-all flex items-center justify-center shadow-lg"
                  >
                    <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </button>
                </div>
              </div>

              {/* Editor Content */}
              <div className="flex-1 overflow-y-auto px-8 py-8">
                <div className="max-w-4xl mx-auto">
                  {/* Title Input */}
                  <input
                    type="text"
                    value={selectedNote.title}
                    onChange={(e) => updateNoteField('title', e.target.value)}
                    placeholder="Untitled Note"
                    className="w-full text-4xl font-bold text-white bg-transparent border-none outline-none placeholder-gray-600 pb-4"
                  />

                  {/* Separator Line */}
                  <div className="mb-6 border-b" style={{ borderColor: 'rgba(227, 208, 149, 0.2)' }}></div>

                  {/* Content Editor */}
                  <div
                    ref={editorRef}
                    id="note-content"
                    contentEditable
                    suppressContentEditableWarning
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        return;
                      }
                    }}
                    onInput={(e) => {
                      saveCursorPosition();
                      updateNoteField('content', e.currentTarget.innerHTML);
                    }}
                    className="w-full min-h-[500px] text-lg text-gray-300 bg-transparent border-none outline-none placeholder-gray-600 resize-none leading-relaxed"
                    style={{
                      whiteSpace: 'pre-wrap',
                      wordWrap: 'break-word'
                    }}
                    dangerouslySetInnerHTML={{ __html: selectedNote.content || '' }}
                  />

                  {/* Footer Info */}
                  <div className="mt-8 pt-4 border-t" style={{ borderColor: 'rgba(255, 255, 255, 0.1)' }}>
                    <p className="text-xs text-gray-500" suppressHydrationWarning>
                      Last edited: {selectedNote.updatedAt ? new Date(selectedNote.updatedAt).toLocaleString() : '-'}
                    </p>
                  </div>
                </div>
              </div>
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <svg className="w-20 h-20 mx-auto mb-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <p className="text-gray-400 text-xl font-medium">Select a note to edit</p>
                <p className="text-gray-500 text-sm mt-2">or create a new one</p>
              </div>
            </div>
          )}
        </div>
      </div>

    </div>
  );
}
