package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Martin4208/Nexus/apps/api-gateway/internal/api"
	"github.com/Martin4208/Nexus/apps/api-gateway/internal/db"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/sqlc-dev/pqtype"
)

// ListGraphs handles GET /workspaces/{workspace_id}/graphs
func (h *Handler) ListGraphs(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID) {
	ctx := r.Context()

	graphs, err := h.queries.ListGraphsByWorkspace(ctx, workspaceId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graphs")
		return
	}

	// Convert to API type
	apiGraphs := make([]api.GraphListItem, len(graphs))
	for i, g := range graphs {
		apiGraphs[i] = convertToAPIGraphListItem(g)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"graphs": apiGraphs,
	})
}

// CreateGraph handles POST /workspaces/{workspace_id}/graphs
func (h *Handler) CreateGraph(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID) {
	ctx := r.Context()

	var reqBody api.CreateGraphJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if reqBody.Title == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title is required")
		return
	}

	// Convert graph_type (*string → sql.NullString)
	var graphType sql.NullString
	if reqBody.GraphType != nil {
		graphType = sql.NullString{String: *reqBody.GraphType, Valid: true}
	} else {
		graphType = sql.NullString{Valid: false}
	}

	graph, err := h.queries.CreateGraph(ctx, db.CreateGraphParams{
		WorkspaceID: workspaceId,
		Title:       reqBody.Title,
		GraphType:   graphType,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create graph")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"graph": convertToAPIGraphFromGraph(graph, nil, nil),
	})
}

// GetGraph handles GET /workspaces/{workspace_id}/graphs/{graph_id}
func (h *Handler) GetGraph(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID, graphId openapi_types.UUID) {
	ctx := r.Context()

	// 1. Get graph metadata
	graph, err := h.queries.GetGraphByID(ctx, graphId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graph")
		return
	}

	// Verify workspace ownership
	if graph.WorkspaceID != workspaceId {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
		return
	}

	// 2. Get nodes
	nodes, err := h.queries.GetGraphNodesByGraphID(ctx, graphId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graph nodes")
		return
	}

	// 3. Get edges
	edges, err := h.queries.GetGraphEdgesByGraphID(ctx, graphId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graph edges")
		return
	}

	// 4. Convert and combine
	apiNodes := make([]api.GraphNode, len(nodes))
	for i, n := range nodes {
		apiNodes[i] = convertToAPIGraphNode(n)
	}

	apiEdges := make([]api.GraphEdge, len(edges))
	for i, e := range edges {
		apiEdges[i] = convertToAPIGraphEdge(e)
	}

	respondJSON(w, http.StatusOK, convertToAPIGraph(graph, apiNodes, apiEdges))
}

// UpdateGraph handles PUT /workspaces/{workspace_id}/graphs/{graph_id}
func (h *Handler) UpdateGraph(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID, graphId openapi_types.UUID) {
	ctx := r.Context()

	var reqBody api.UpdateGraphJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Verify graph exists and belongs to workspace
	existing, err := h.queries.GetGraphByID(ctx, graphId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graph")
		return
	}

	if existing.WorkspaceID != workspaceId {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
		return
	}

	// Convert layout_config to pqtype.NullRawMessage
	var layoutConfig pqtype.NullRawMessage
	if reqBody.LayoutConfig != nil {
		layoutBytes, _ := json.Marshal(reqBody.LayoutConfig)
		layoutConfig = pqtype.NullRawMessage{RawMessage: layoutBytes, Valid: true}
	}

	graph, err := h.queries.UpdateGraph(ctx, db.UpdateGraphParams{
		ID:           graphId,
		Title:        sql.NullString{String: stringValue(reqBody.Title), Valid: reqBody.Title != nil},
		LayoutConfig: layoutConfig,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update graph")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"graph": convertToAPIGraphFromGraph(graph, nil, nil),
	})
}

// DeleteGraph handles DELETE /workspaces/{workspace_id}/graphs/{graph_id}
func (h *Handler) DeleteGraph(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID, graphId openapi_types.UUID) {
	ctx := r.Context()

	// Verify ownership
	graph, err := h.queries.GetGraphByID(ctx, graphId)
	if err != nil {
		if err == sql.ErrNoRows {
			respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch graph")
		return
	}

	if graph.WorkspaceID != workspaceId {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Graph not found")
		return
	}

	if err := h.queries.DeleteGraph(ctx, graphId); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete graph")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetGraphNodes handles GET /workspaces/{workspace_id}/graphs/{graph_id}/nodes
func (h *Handler) GetGraphNodes(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID, graphId openapi_types.UUID) {
	ctx := r.Context()

	nodes, err := h.queries.GetGraphNodesByGraphID(ctx, graphId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch nodes")
		return
	}

	apiNodes := make([]api.GraphNode, len(nodes))
	for i, n := range nodes {
		apiNodes[i] = convertToAPIGraphNode(n)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"nodes": apiNodes,
	})
}

// GetGraphEdges handles GET /workspaces/{workspace_id}/graphs/{graph_id}/edges
func (h *Handler) GetGraphEdges(w http.ResponseWriter, r *http.Request, workspaceId openapi_types.UUID, graphId openapi_types.UUID) {
	ctx := r.Context()

	edges, err := h.queries.GetGraphEdgesByGraphID(ctx, graphId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch edges")
		return
	}

	apiEdges := make([]api.GraphEdge, len(edges))
	for i, e := range edges {
		apiEdges[i] = convertToAPIGraphEdge(e)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"edges": apiEdges,
	})
}

// ========================================
// Conversion Helpers
// ========================================

func convertToAPIGraphListItem(g db.ListGraphsByWorkspaceRow) api.GraphListItem {
	return api.GraphListItem{
		Id:          g.ID,
		WorkspaceId: g.WorkspaceID,
		Title:       g.Title,
		GraphType:   nullStringToPtr(g.GraphType),
		CreatedAt:   g.CreatedAt, // time.Time (ポインタなし)
		UpdatedAt:   g.UpdatedAt, // time.Time (ポインタなし)
	}
}

func convertToAPIGraph(g db.GetGraphByIDRow, nodes []api.GraphNode, edges []api.GraphEdge) api.Graph {
	var layoutConfig *map[string]interface{}
	if g.LayoutConfig.Valid {
		var config map[string]interface{}
		if err := json.Unmarshal(g.LayoutConfig.RawMessage, &config); err == nil {
			layoutConfig = &config
		}
	}

	return api.Graph{
		Id:           g.ID,
		WorkspaceId:  g.WorkspaceID,
		Title:        g.Title,
		GraphType:    nullStringToPtr(g.GraphType),
		LayoutConfig: layoutConfig,
		Nodes:        nodes,
		Edges:        edges,
		CreatedAt:    g.CreatedAt, // time.Time (ポインタなし)
		UpdatedAt:    g.UpdatedAt, // time.Time (ポインタなし)
	}
}

func convertToAPIGraphNode(n db.GraphNode) api.GraphNode {
	// Position の変換（float64 → float32）
	var position *struct {
		X *float32 `json:"x,omitempty"`
		Y *float32 `json:"y,omitempty"`
		Z *float32 `json:"z"`
	}

	if n.Position.Valid {
		var pos struct {
			X *float64 `json:"x,omitempty"`
			Y *float64 `json:"y,omitempty"`
			Z *float64 `json:"z,omitempty"`
		}
		if err := json.Unmarshal(n.Position.RawMessage, &pos); err == nil {
			position = &struct {
				X *float32 `json:"x,omitempty"`
				Y *float32 `json:"y,omitempty"`
				Z *float32 `json:"z"`
			}{
				X: float64ToFloat32Ptr(pos.X),
				Y: float64ToFloat32Ptr(pos.Y),
				Z: float64ToFloat32Ptr(pos.Z),
			}
		}
	}

	var style *map[string]interface{}
	if n.Style.Valid {
		var s map[string]interface{}
		if err := json.Unmarshal(n.Style.RawMessage, &s); err == nil {
			style = &s
		}
	}

	var metadata *map[string]interface{}
	if n.Metadata.Valid {
		var m map[string]interface{}
		if err := json.Unmarshal(n.Metadata.RawMessage, &m); err == nil {
			metadata = &m
		}
	}

	return api.GraphNode{
		Id:         n.ID,
		GraphId:    n.GraphID,
		Label:      n.Label,
		NodeType:   n.NodeType,
		SourceType: nullStringToPtr(n.SourceType),
		SourceId:   nullUUIDToPtr(n.SourceID),
		Position:   position,
		Style:      style,
		Metadata:   metadata,
		CreatedAt:  timeToPtr(n.CreatedAt), // time.Time → *time.Time
	}
}

func convertToAPIGraphEdge(e db.GraphEdge) api.GraphEdge {
	var style *map[string]interface{}
	if e.Style.Valid {
		var s map[string]interface{}
		if err := json.Unmarshal(e.Style.RawMessage, &s); err == nil {
			style = &s
		}
	}

	var metadata *map[string]interface{}
	if e.Metadata.Valid {
		var m map[string]interface{}
		if err := json.Unmarshal(e.Metadata.RawMessage, &m); err == nil {
			metadata = &m
		}
	}

	isDirected := e.IsDirected

	return api.GraphEdge{
		Id:         e.ID,
		GraphId:    e.GraphID,
		FromNodeId: e.FromNodeID,
		ToNodeId:   e.ToNodeID,
		EdgeType:   e.EdgeType,
		IsDirected: &isDirected,                       // bool → *bool
		Weight:     nullFloat64ToFloat32Ptr(e.Weight), // sql.NullFloat64 → *float32
		Confidence: nullFloat64ToFloat32Ptr(e.Confidence),
		Style:      style,
		Metadata:   metadata,
		CreatedAt:  timeToPtr(e.CreatedAt), // time.Time → *time.Time
		UpdatedAt:  timeToPtr(e.UpdatedAt), // time.Time → *time.Time
	}
}

// ========================================
// Null Type Helpers
// ========================================

func nullStringToPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullUUIDToPtr(nu uuid.NullUUID) *openapi_types.UUID {
	if nu.Valid {
		return &nu.UUID
	}
	return nil
}

func nullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}

// sql.NullFloat64 → *float32
func nullFloat64ToFloat32Ptr(nf sql.NullFloat64) *float32 {
	if nf.Valid {
		f32 := float32(nf.Float64)
		return &f32
	}
	return nil
}

// *float64 → *float32
func float64ToFloat32Ptr(f *float64) *float32 {
	if f == nil {
		return nil
	}
	f32 := float32(*f)
	return &f32
}

// time.Time → *time.Time
func timeToPtr(t time.Time) *time.Time {
	return &t
}

func stringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// convertToAPIGraph のオーバーロード（db.Graph用）
func convertToAPIGraphFromGraph(g db.Graph, nodes []api.GraphNode, edges []api.GraphEdge) api.Graph {
	var layoutConfig *map[string]interface{}
	if g.LayoutConfig.Valid {
		var config map[string]interface{}
		if err := json.Unmarshal(g.LayoutConfig.RawMessage, &config); err == nil {
			layoutConfig = &config
		}
	}

	return api.Graph{
		Id:           g.ID,
		WorkspaceId:  g.WorkspaceID,
		Title:        g.Title,
		GraphType:    nullStringToPtr(g.GraphType),
		LayoutConfig: layoutConfig,
		Nodes:        nodes,
		Edges:        edges,
		CreatedAt:    g.CreatedAt,
		UpdatedAt:    g.UpdatedAt,
	}
}
