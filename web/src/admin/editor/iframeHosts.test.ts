import { describe, expect, it } from 'vitest'
import { isAllowedIframe, toEmbedUrl } from './iframeHosts'

describe('isAllowedIframe', () => {
  it.each([
    ['https://www.youtube.com/embed/abc', true],
    ['//www.youtube.com/embed/abc', true],
    ['https://players.brightcove.net/1/x/index.html?videoId=2', true],
    ['https://drive.google.com/file/d/x/preview', true],
    ['https://www.youtube.com.evil.com/embed', false],
    ['https://evil.example/embed', false],
    ['/embed/abc', false],
    ['javascript:alert(1)', false],
    ['', false],
  ])('%s → %s', (src, want) => {
    expect(isAllowedIframe(src)).toBe(want)
  })
})

describe('toEmbedUrl', () => {
  it.each([
    ['https://www.youtube.com/watch?v=V4LhnGIPZt4&t=10s', 'https://www.youtube.com/embed/V4LhnGIPZt4'],
    ['https://youtu.be/V4LhnGIPZt4', 'https://www.youtube.com/embed/V4LhnGIPZt4'],
    ['https://www.youtube.com/shorts/abc123', 'https://www.youtube.com/embed/abc123'],
    ['  https://www.youtube.com/embed/x  ', 'https://www.youtube.com/embed/x'],
    ['<iframe width="560" src="https://www.youtube.com/embed/x?si=1" frameborder="0"></iframe>', 'https://www.youtube.com/embed/x?si=1'],
    ['https://drive.google.com/file/d/x/preview', 'https://drive.google.com/file/d/x/preview'],
  ])('%s', (input, want) => {
    expect(toEmbedUrl(input)).toBe(want)
  })
})
