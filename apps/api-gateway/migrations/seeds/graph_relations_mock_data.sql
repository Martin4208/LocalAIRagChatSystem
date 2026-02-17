-- Graph Relations Mock Data
-- Workspace: 967e7e13-dc87-420c-9994-64d99c6b9e35

-- まず既存EntityのIDを定数として定義（見やすくするため）
DO $$
DECLARE
    ws_id UUID := '967e7e13-dc87-420c-9994-64d99c6b9e35';
    sato_id UUID;
    tanaka_id UUID;
    suzuki_id UUID;
    nexus_id UUID;
    ai_solutions_id UUID;
BEGIN
    -- EntityのIDを取得
    SELECT id INTO sato_id FROM graph_entities WHERE label = '佐藤健' AND workspace_id = ws_id;
    SELECT id INTO tanaka_id FROM graph_entities WHERE label = '田中花子' AND workspace_id = ws_id;
    SELECT id INTO suzuki_id FROM graph_entities WHERE label = '鈴木一郎' AND workspace_id = ws_id;
    SELECT id INTO nexus_id FROM graph_entities WHERE label = '株式会社Nexus' AND workspace_id = ws_id;
    SELECT id INTO ai_solutions_id FROM graph_entities WHERE label = '株式会社AI Solutions' AND workspace_id = ws_id;
    
    -- Relationを挿入
    INSERT INTO graph_relations (workspace_id, source_entity_id, target_entity_id, type, is_directed, weight, metadata)
    VALUES
        (ws_id, sato_id, nexus_id, 'works_for', TRUE, 1.0, '{"role": "CTO"}'),
        (ws_id, tanaka_id, nexus_id, 'works_for', TRUE, 1.0, '{"role": "CEO"}'),
        (ws_id, suzuki_id, nexus_id, 'works_for', TRUE, 1.0, '{"role": "Senior Engineer"}'),
        (ws_id, suzuki_id, sato_id, 'reports_to', TRUE, 1.0, '{}'),
        (ws_id, nexus_id, ai_solutions_id, 'partnered_with', FALSE, 1.0, '{"type": "Strategic"}');
END $$;