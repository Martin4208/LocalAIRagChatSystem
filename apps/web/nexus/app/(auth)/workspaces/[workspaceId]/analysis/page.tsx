'use client';

import { useState, useRef, useEffect } from 'react';
import { useParams } from 'next/navigation';

import { Card, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function Analysis({ documents }) {
    const [analysis, setAnalysis] = useState('');

    if (!documents || documents.length === 0) {
        return (
            <Card>
                <div>
                    <p>No documents to analyze</p>
                    <p>Upload documents to generate insights</p>
                </div>
            </Card>
        );
    }

    return (
        <Card>
            <CardHeader>
                <div>
                    AI Analysis
                </div>
                <Button
                    type="submit"
                    disabled={isLoading}
                >
                    {isLoading ? (
                        <>
                            Analyzing...
                        </>
                    ) : (
                        <>
                            {analysis ? 'Refresh' : 'Generate'} analysis
                        </>
                    )}
                </Button>
            </CardHeader>

            <div>
                This is where the analysis will come
            </div>
        </Card>
    );
}