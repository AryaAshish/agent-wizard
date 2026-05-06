#!/usr/bin/env node
'use strict';

const fs = require('fs');
const path = require('path');
const os = require('os');
const crypto = require('crypto');
const { spawnSync } = require('child_process');

const REPO_DEFAULT = 'AryaAshish/agent-wizard';

async function fetchBuffer(url, headers = {}) {
  const res = await fetch(url, { redirect: 'follow', headers });
  if (!res.ok) {
    const t = await res.text();
    throw new Error(`Fetch ${url}: ${res.status} ${t.slice(0, 200)}`);
  }
  return Buffer.from(await res.arrayBuffer());
}

async function resolveVersion(repo) {
  const fromEnv = process.env.AGENT_WIZARD_VERSION;
  if (fromEnv && fromEnv !== 'latest') {
    return String(fromEnv).replace(/^v/, '');
  }
  const res = await fetch(`https://github.com/${repo}/releases/latest`, { redirect: 'follow' });
  const finalUrl = res.url || '';
  const segment = (finalUrl.split('/').pop() || '').trim();
  const stripped = segment.replace(/^v/, '');
  const looksReasonable =
    stripped.length > 0 &&
    stripped.length < 128 &&
    !/[/?#]|:/.test(stripped);

  if (looksReasonable) {
    return stripped;
  }

  const gh = await fetch(`https://api.github.com/repos/${repo}/releases/latest`, {
    headers: { Accept: 'application/vnd.github+json' },
  });
  const j = await gh.json().catch(() => ({}));
  const tag =
    gh.ok && typeof j.tag_name === 'string' && j.tag_name.length > 0 ? j.tag_name : null;

  if (!tag) {
    throw new Error(
      `Could not resolve latest release for ${repo}; set AGENT_WIZARD_VERSION to a tag (without v prefix is ok).`,
    );
  }

  return tag.replace(/^v/, '');
}

function platformAsset() {
  const p = process.platform;
  if (p === 'darwin') return { os: 'darwin', archiveExt: '.tar.gz', binName: 'agent-wizard' };
  if (p === 'linux') return { os: 'linux', archiveExt: '.tar.gz', binName: 'agent-wizard' };
  if (p === 'win32') return { os: 'windows', archiveExt: '.zip', binName: 'agent-wizard.exe' };
  throw new Error(`Unsupported platform: ${p}`);
}

function normalizedArch() {
  const a = os.arch();
  if (a === 'x64') return 'amd64';
  if (a === 'arm64') return 'arm64';
  throw new Error(`Unsupported arch: ${a}`);
}

function sha256(buf) {
  return crypto.createHash('sha256').update(buf).digest('hex');
}

function parseChecksums(text, assetName) {
  for (const line of text.split('\n')) {
    const m = line.match(/^([a-fA-F0-9]{64})\s+(.+)$/);
    if (m && m[2].trim() === assetName) return m[1].toLowerCase();
  }
  return null;
}

function extract(archivePath, destDir, archiveExt) {
  fs.mkdirSync(destDir, { recursive: true });
  if (archiveExt === '.tar.gz') {
    const r = spawnSync('tar', ['-xzf', archivePath, '-C', destDir], { stdio: 'inherit' });
    if (r.status !== 0) throw new Error('Failed to extract release archive');
    return;
  }
  const u = spawnSync('unzip', ['-q', '-o', archivePath, '-d', destDir], {
    stdio: 'pipe',
    encoding: 'utf8',
  });
  if (u.status === 0) return;

  const destEsc = destDir.replace(/'/g, "''");
  const arcEsc = archivePath.replace(/'/g, "''");
  const psCmd = `Expand-Archive -LiteralPath '${arcEsc}' -DestinationPath '${destEsc}' -Force`;
  const ps = spawnSync('powershell.exe', ['-NoProfile', '-NonInteractive', '-Command', psCmd], {
    stdio: 'inherit',
    windowsHide: true,
  });
  if (ps.status !== 0) {
    throw new Error(
      `Failed to extract Windows archive (needs unzip or PowerShell Expand-Archive). ${u.stderr || u.stdout || ''}`,
    );
  }
}

async function ensureBinary(repo, cacheRoot) {
  const { os: goos, archiveExt, binName } = platformAsset();
  const goarch = normalizedArch();
  const version = await resolveVersion(repo);
  const archiveName = `agent-wizard_${version}_${goos}_${goarch}${archiveExt}`;
  const versionDir = path.join(cacheRoot, `v${version}`);
  const binPath = path.join(versionDir, binName);

  if (fs.existsSync(binPath)) {
    try {
      if (goos !== 'windows') fs.chmodSync(binPath, 0o755);
    } catch (_) {
      /* best effort */
    }
    return binPath;
  }

  const baseUrl = `https://github.com/${repo}/releases/download/v${version}`;
  const [archiveBuf, checksumText] = await Promise.all([
    fetchBuffer(`${baseUrl}/${archiveName}`),
    fetchBuffer(`${baseUrl}/checksums.txt`).then((b) => b.toString('utf8')),
  ]);

  const expected = parseChecksums(checksumText, archiveName);
  if (!expected) throw new Error(`Checksum not found for ${archiveName}`);

  const actual = sha256(archiveBuf).toLowerCase();
  if (actual !== expected) throw new Error('Checksum verification failed for downloaded archive');

  const tmpDir = path.join(cacheRoot, `.tmp-${process.pid}`);
  fs.mkdirSync(tmpDir, { recursive: true });
  try {
    const arcFile = path.join(tmpDir, archiveName);
    fs.writeFileSync(arcFile, archiveBuf);
    extract(arcFile, versionDir, archiveExt);

    try {
      if (goos !== 'windows') fs.chmodSync(binPath, 0o755);
    } catch (_) {
      /* best effort */
    }

    if (!fs.existsSync(binPath)) {
      throw new Error(`Binary missing after extract: ${binPath}`);
    }
    return binPath;
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

async function main() {
  const repo = process.env.AGENT_WIZARD_REPO || REPO_DEFAULT;
  const cacheRoot =
    process.env.AGENT_WIZARD_CACHE_DIR ||
    path.join(os.homedir(), '.cache', 'agent-wizard', 'npm');

  fs.mkdirSync(cacheRoot, { recursive: true });

  const binPath = await ensureBinary(repo, cacheRoot);

  const r = spawnSync(binPath, process.argv.slice(2), {
    stdio: 'inherit',
    env: process.env,
    windowsHide: true,
    shell: false,
  });

  process.exit(Number.isInteger(r.status) ? r.status : 1);
}

main().catch((e) => {
  console.error(e instanceof Error ? e.message : String(e));
  process.exit(1);
});
