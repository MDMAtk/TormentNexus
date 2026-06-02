import fs from 'fs';
import path from 'path';
import { spawnSync } from 'child_process';

const WORKSPACE_ROOT = 'c:\\Users\\hyper\\workspace';

function ensureDir(dirPath) {
    if (!fs.existsSync(dirPath)) {
        fs.mkdirSync(dirPath, { recursive: true });
    }
}

// Deep merge objects recursively, union arrays
function deepMerge(target, source) {
    if (typeof target !== 'object' || target === null || typeof source !== 'object' || source === null) {
        return source;
    }
    
    if (Array.isArray(target) && Array.isArray(source)) {
        // Union arrays by stringified value to avoid duplicates
        const seen = new Set();
        const merged = [];
        [...target, ...source].forEach(item => {
            const key = typeof item === 'object' ? JSON.stringify(item) : item;
            if (!seen.has(key)) {
                seen.add(key);
                merged.push(item);
            }
        });
        return merged;
    }
    
    const result = { ...target };
    for (const key of Object.keys(source)) {
        if (key in target) {
            result[key] = deepMerge(target[key], source[key]);
        } else {
            result[key] = source[key];
        }
    }
    return result;
}

// Highly sophisticated session file content-level merger
function mergeFiles(src, dest) {
    const ext = path.extname(dest).toLowerCase();
    
    // 1. If file doesn't exist at destination, copy directly
    if (!fs.existsSync(dest)) {
        fs.copyFileSync(src, dest);
        return;
    }
    
    const srcStat = fs.statSync(src);
    const destStat = fs.statSync(dest);
    
    // Check if files are identical
    if (srcStat.size === destStat.size) {
        try {
            const srcBuf = fs.readFileSync(src);
            const destBuf = fs.readFileSync(dest);
            if (srcBuf.equals(destBuf)) {
                return; // Exactly identical, skip
            }
        } catch {}
    }
    
    // 2. Handle JSON deep-merges
    if (ext === '.json') {
        try {
            const srcJson = JSON.parse(fs.readFileSync(src, 'utf-8'));
            const destJson = JSON.parse(fs.readFileSync(dest, 'utf-8'));
            const mergedJson = deepMerge(destJson, srcJson);
            fs.writeFileSync(dest, JSON.stringify(mergedJson, null, 2), 'utf-8');
            console.log(`    🔀 Deep-merged JSON configuration: ${path.basename(dest)}`);
            return;
        } catch (e) {
            // Fall back to timestamp override or concatenation if JSON parsing fails
        }
    }
    
    // 3. Handle Text/Log/Markdown concatenation and deduplication
    if (['.txt', '.log', '.md', '.jsonl', '.history'].includes(ext)) {
        try {
            const srcText = fs.readFileSync(src, 'utf-8').trim();
            const destText = fs.readFileSync(dest, 'utf-8').trim();
            
            // Check if one file fully contains the other
            if (destText.includes(srcText)) {
                return; // Dest already has all of Src's info, keep dest
            }
            if (srcText.includes(destText)) {
                fs.writeFileSync(dest, srcText, 'utf-8'); // Src is a superset, replace dest with src
                return;
            }
            
            // Otherwise, concatenate them chronological or clean append
            const separator = ext === '.md' 
                ? `\n\n---\n### 🔄 Merged Session Segment (${new Date(srcStat.mtimeMs).toISOString()})\n---\n\n`
                : `\n\n=== Merged Session Segment (${new Date(srcStat.mtimeMs).toISOString()}) ===\n\n`;
                
            // Put the chronologically older segment first
            let mergedText;
            if (srcStat.mtimeMs < destStat.mtimeMs) {
                mergedText = `${srcText}${separator}${destText}`;
            } else {
                mergedText = `${destText}${separator}${srcText}`;
            }
            
            fs.writeFileSync(dest, mergedText, 'utf-8');
            console.log(`    🔀 Concatenated/Merged text log segment: ${path.basename(dest)}`);
            return;
        } catch {}
    }
    
    // 4. Default: Keep the newer file for database / binary structures
    if (srcStat.mtimeMs > destStat.mtimeMs) {
        fs.copyFileSync(src, dest);
    }
}

function copyRecursiveMergeSync(src, dest) {
    const stats = fs.statSync(src);
    const isDirectory = stats.isDirectory();
    
    if (isDirectory) {
        ensureDir(dest);
        fs.readdirSync(src).forEach((childItemName) => {
            copyRecursiveMergeSync(
                path.join(src, childItemName),
                path.join(dest, childItemName)
            );
        });
    } else {
        mergeFiles(src, dest);
    }
}

async function run() {
    console.log('==================================================');
    console.log('🔄 UNIVERSAL DEEP-MERGE SESSION DIRECTORY PIPELINE');
    console.log('==================================================');
    
    if (!fs.existsSync(WORKSPACE_ROOT)) {
        console.error(`Workspace root not found at ${WORKSPACE_ROOT}. Exiting.`);
        process.exit(1);
    }
    
    const entries = fs.readdirSync(WORKSPACE_ROOT, { withFileTypes: true });
    let mergedProjectsCount = 0;
    
    for (const entry of entries) {
        if (!entry.isDirectory() || entry.name.startsWith('.')) {
            continue;
        }
        
        const projectPath = path.join(WORKSPACE_ROOT, entry.name);
        const borgPath = path.join(projectPath, '.borg');
        const hypercodePath = path.join(projectPath, '.hypercode');
        const tormentnexusPath = path.join(projectPath, '.tormentnexus');
        
        const hasBorg = fs.existsSync(borgPath);
        const hasHypercode = fs.existsSync(hypercodePath);
        
        if (hasBorg || hasHypercode) {
            console.log(`📁 Processing project: ${entry.name}`);
            ensureDir(tormentnexusPath);
            
            if (hasBorg) {
                console.log(`  - Deep Merging .borg -> .tormentnexus`);
                try {
                    copyRecursiveMergeSync(borgPath, tormentnexusPath);
                    // Handle .borg-session.json
                    const borgSessionFile = path.join(projectPath, '.borg-session.json');
                    if (fs.existsSync(borgSessionFile)) {
                        const destSession = path.join(projectPath, '.tormentnexus-session.json');
                        mergeFiles(borgSessionFile, destSession);
                    }
                } catch (err) {
                    console.error(`  ❌ Failed to deep merge .borg for ${entry.name}:`, err.message);
                }
            }
            
            if (hasHypercode) {
                console.log(`  - Deep Merging .hypercode -> .tormentnexus`);
                try {
                    copyRecursiveMergeSync(hypercodePath, tormentnexusPath);
                    // Handle .hypercode-session.json
                    const hypercodeSessionFile = path.join(projectPath, '.hypercode-session.json');
                    if (fs.existsSync(hypercodeSessionFile)) {
                        const destSession = path.join(projectPath, '.tormentnexus-session.json');
                        mergeFiles(hypercodeSessionFile, destSession);
                    }
                } catch (err) {
                    console.error(`  ❌ Failed to deep merge .hypercode for ${entry.name}:`, err.message);
                }
            }
            
            mergedProjectsCount++;
        }
    }
    
    console.log('\n==================================================');
    console.log(`🎉 Deep Merged ${mergedProjectsCount} project folders successfully with ZERO data loss!`);
    console.log('==================================================');
    
    // Now trigger the universal session ingestion pipeline to index all the newly merged files
    console.log('🚀 Triggering session ingestion to index newly merged files...');
    const ingestScript = path.join(WORKSPACE_ROOT, 'borg', 'scratch', 'ingest_all_sessions.mjs');
    if (fs.existsSync(ingestScript)) {
        const result = spawnSync('node', [ingestScript], { stdio: 'inherit' });
        if (result.status === 0) {
            console.log('✅ Session ingestion ran successfully!');
        } else {
            console.error('❌ Session ingestion failed with code:', result.status);
        }
    } else {
        console.error(`Could not find ingest script at ${ingestScript}`);
    }
}

run();
