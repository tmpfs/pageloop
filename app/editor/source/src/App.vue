<template>
  <main>
    <app-alert></app-alert>
    <app-notify></app-notify>
    <app-header></app-header>
    <app-main></app-main>
    <app-footer></app-footer>
  </main>
</template>

<script>
import AppAlert from '@/components/Alert'
import AppNotify from '@/components/Notify'
import AppHeader from '@/components/Header'
import AppMain from '@/components/Main'
import AppFooter from '@/components/Footer'

window.addEventListener('beforeunload', (e) => {
  e.preventDefault()
  return false
})

export default {
  name: 'app',
  components: {AppAlert, AppNotify, AppHeader, AppMain, AppFooter}
}
</script>

<style>
  /* Variables */
  :root {
    --base03-color: #002b36;
    --base02-color: #073642;
    --base01-color: #586e75;
    --base00-color: #657b83;
    --base0-color: #839496;
    --base1-color: #93a1a1;
    --base2-color: #eee8d5;
    --base3-color: #fdf6e3;
    --yellow-color: #b58900;
    --orange-color: #cb4b16;
    --red-color: #dc322f;
    --magenta-color: #d33682;
    --violet-color: #6c71c4;
    --blue-color: #268bd2;
    --cyan-color: #2aa198;
    --green-color: #859900;

    --background: #073642;
    --color: #93a1a1;
    --border-color: #002b36;
    --toolbar-link: #839496;

    --drag-size: 6px;
  }

  /* Fonts */

  @font-face {
    font-family: 'SourceSerif';
    src: url('/assets/font/sourceserifpro-regular-webfont.woff2') format('woff2'),
         url('/assets/font/sourceserifpro-regular-webfont.woff') format('woff');
    font-style: normal;
    font-weight: normal;
    text-rendering: optimizeLegibility;
  }

  @font-face {
    font-family: 'SourceSans';
    src: url('/assets/font/sourcesanspro-regular-webfont.woff2') format('woff2'),
         url('/assets/font/sourcesanspro-regular-webfont.woff') format('woff');
    font-style: normal;
    font-weight: normal;
    text-rendering: optimizeLegibility;
  }

  /* Document */

  html {
    box-sizing: border-box;
    font-size: 62.5%;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  *, *:after, *:before {
    box-sizing: inherit
  }

  html, body {
    width: 100%;
    height: 100%;
    min-height: 100%;
    margin: 0;
    padding: 0;
  }

  body {
    font-size: 1.8rem;
    line-height: 2.2rem;
    font-family: 'SourceSans';
    background: var(--background);
    color: var(--color);
    overflow: hidden;
  }

  main {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    height: 100%;
    min-height: 100%;
  }

  main > .content-main {
    flex: 1 1;
    height: calc(100% - 6rem);
  }

  .content {
    position: relative;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    overflow: hidden;
    height: 100%;
    min-height: 100%;
  }

  .scroll {
    overflow-x: hidden;
    overflow-y: auto;
    height: calc(100% - 4.6rem);
  }

  /* Links */

  a, a:visited, a:hover {
    color: var(--blue-color);
    text-decoration: none;
    transition: all 0.5s ease-out;
    cursor: pointer;
  }

  p a{
    display: inline-block;
    line-height: 1.8rem;
    border-bottom: 1px solid transparent;
  }

  p a:hover {
    border-bottom: 1px solid currentColor;
  }

  a.close {
    display: inline-block;
    color: var(--color);
    width: 3.2rem;
    height: 3.2rem;
    padding: 0.7rem 0.5rem 0.5rem 0.5rem;
  }

  a.close::before {
    content: '🞩';
    padding-left: 0.5rem;
  }

  a.close:hover {
    color: var(--base3-color);
    background: var(--base03-color);
  }

  /* Headings */

  h1, h2, h3, h4, h5, h6 {
    font-weight: normal;
    margin: 0;
    padding: 0;
  }

  h1 {
    font-family: SourceSans;
    text-align: center;
    margin: 3rem 0 0 0;
    font-size: 2.6rem;
    vertical-align: top;
    border-bottom: 4px double currentColor;
    padding-bottom: 3rem;
    text-transform: uppercase;
  }

  h1 > span:first-child {
    font-size: 2.4rem;
  }

  /* Forms & Buttons */

  form, fieldset {
    margin: 0;
    padding: 0;
    border: 0;
  }

  form > p {
    margin-bottom: 0.5rem;
  }

  input[type="text"], input[type="submit"], input[type="reset"], button {
    font-size: 1.6rem;
    font-family: SourceSans;
    padding: 0.6rem;
    border-radius: 4px;
    background: var(--base2-color);
    color: var(--base03-color);
  }

  input[type="text"] {
    width: 100%;
  }

  button, input[type="submit"], input[type="reset"] {
    cursor: pointer;
    border: 2px solid var(--base1-color);
    border-style: none;
    background: var(--base03-color);
    color: var(--base2-color);
    min-width: 6rem;
  }

  button:hover {
    color: var(--base3-color);
  }

  button.sml {
    font-size: 1.4rem;
    padding: 0.3rem;
  }

  button.primary, input[type="submit"].primary {
    background: var(--blue-color);
    color: var(--base3-color);
  }

  input[type="submit"], input[type="reset"] {
    margin-top: 1rem;
  }

  /* Rhythm */

  section {
    margin-bottom: 3rem;
  }

  p, section > ul.compact-list {
    margin: 1rem 0;
  }

  section.document {
    font-family: SourceSerif;
    text-align: left;
    display: inline-block;
    margin: auto 0;
    min-width: 480px;
    max-width: 640px;
  }

  /* Misc Elements */

  iframe {
    display: block;
    border: 0;
    margin: 0;
    width: 100%;
  }

  /* Utilities */

  .hidden {
    display: none !important;
  }

  .disabled {
    opacity: 0.6;
    pointer-events: none;
  }

  .small, small {
    font-size: 1.4rem;
  }

  .compact-list {
    list-style-type: none;
    margin: 0;
    padding: 0;
  }

  .clearfix {
    clear: both;
  }

  .selected {
    pointer-events: none;
  }

  .center {
    text-align: center;
  }

  .form-actions {
    text-align: right;
    margin-top: 1rem;
  }

  .form-actions > * {
    margin-left: 0.5rem;
  }

  .danger {
    background: var(--red-color);
  }

  .panel {
    padding: 2rem;
  }
</style>
