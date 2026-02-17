### DB設計 ###

### CREATE TABLE 

- workspace
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    settings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);


- CREATE TABLE directories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES directories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    
    UNIQUE(workspace_id, parent_id, name)
);


- CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    directory_id UUID REFERENCES directories(id) ON DELETE SET NULL,
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    tags TEXT[],
    metadata JSONB,
    status TEXT NOT NULL DEFAULT 'uploaded',
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_documents_workspace ON documents(workspace_id);
CREATE INDEX idx_documents_directory ON documents(directory_id);
CREATE INDEX idx_documents_tags ON documents USING GIN(tags);
CREATE INDEX idx_documents_status ON documents(status);


- CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    sha256_hash TEXT NOT NULL UNIQUE,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    original_filename TEXT,
    minio_bucket TEXT NOT NULL,
    minio_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_accessed_at TIMESTAMPTZ
);
CREATE INDEX idx_files_sha256 ON files(sha256_hash);

- CREATE TABLE document_chunks(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    chunk_index INTEGER NOT NULL,
    content TEXT NOT NULL,
    qdrant_point_id UUID,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(document_id, chunk_index)
);
CREATE INDEX idx_chunks_document ON document_chunks(document_id);

- CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    filter_config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

- CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    message_index INTEGER NOT NULL,
    document_refs JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(chat_id, message_index)
);
CREATE INDEX idx_messages_chat ON chat_messages(chat_id, message_index);

- CREATE TABLE graphs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    graph_type TEXT,
    layout_config JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

- CREATE TABLE graph_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    graph_id UUID  NOT NULL REFERENCES graphs(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    node_type TEXT NOT NULL,
    source_type TEXT,
    source_id UUID,
    position JSONB,
    style JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_source FOREIGN KEY (source_id) REFERENCES document_chunks(id) ON DELETE SET NULL
);
CREATE INDEX idx_graph_nodes_graph ON graph_nodes(graph_id);
CREATE INDEX idx_graph_nodes_source ON graph_nodes(source_id);

- CREATE TABLE graph_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    graph_id UUID NOT NULL REFERENCES graphs(id) ON DELETE CASCADE,
    from_node_id UUID NOT NULL REFERENCES graph_nodes(id) ON DELETE CASCADE,
    to_node_id UUID NOT NULL REFERENCES graph_nodes(id) ON DELETE CASCADE, 
    edge_type TEXT NOT NULL,
    is_directed BOOLEAN NOT NULL DEFAULT true,
    weight FLOAT,
    confidence FLOAT,
    style JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(graph_id, from_node_id, to_node_id, edge_type)
);
CREATE INDEX idx_graph_edges_graph ON graph_edges(graph_id);
CREATE INDEX idx_graph_edges_from ON graph_edges(from_node_id);
CREATE INDEX idx_graph_edges_to ON graph_edges(to_node_id);

- CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    analysis_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    config JSONB,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_analyses_workspace ON analyses(workspace_id);
CREATE INDEX idx_analyses_status ON analyses(status);

- CREATE TABLE analysis_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    result_type TEXT NOT NULL,
    content JSONB NOT NULL,
    image_url TEXT,
    minio_bucket TEXT,
    minio_key TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_analysis_results_analysis ON analysis_results(analysis_id);
CREATE INDEX idx_analysis_results_type ON analysis_results(result_type);

- CREATE TABLE canvases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    settings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

- CREATE TABLE canvas_elements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    canvas_id UUID NOT NULL REFERENCES canvases(id) ON DELETE CASCADE,
    element_type TEXT NOT NULL,
    position JSONB NOT NULL,
    z_index INTEGER NOT NULL DEFAULT 0,
    content JSONB NOT NULL,
    style JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_canvas_elements_canvas ON canvas_elements(canvas_id);
CREATE INDEX idx_canvas_elements_type ON canvas_elements(element_type);