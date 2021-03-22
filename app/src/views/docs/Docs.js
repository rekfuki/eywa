import { useState, useEffect } from 'react';
import { useLocation, useHistory } from 'react-router-dom';
import Markdown from 'react-markdown';
import matter from 'gray-matter';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import dracula from 'react-syntax-highlighter/dist/cjs/styles/prism/dracula';
import { Container, useTheme, makeStyles } from '@material-ui/core';
import styled from 'styled-components';

// const MarkdownWrapper = styled('div')((theme) => ({
const MarkdownWrapper = styled('div')`
  ${({ theme }) => `
  color: ${theme.palette.text.primary};
  font-family: ${theme.typography.fontFamily};
  & blockquote {
    border-left: 4px solid ${theme.palette.text.secondary};
    margin-bottom: ${theme.spacing(2)}px;
    padding-bottom: ${theme.spacing(1)}px;
    padding-left: ${theme.spacing(2)}px;
    padding-top: ${theme.spacing(1)}px;
    & > p {
      color: ${theme.palette.text.secondary};
      margin-bottom: 0px;
    }
  }
  & code {
    color: #01ab56;
    font-family: Inconsolata; Monaco; Consolas; \'Courier New\'; Courier; monospace;
    font-size: 14px;
    padding-left: 2px;
    padding-right: 2px;
  }
  & h1 {
    font-size: 35px;
    font-weight: 500;
    letter-spacing: -0.24px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(6)}px;
  }
  & h2 {
    font-size: 29px;
    font-weight: 500;
    letter-spacing: -0.24px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(6)}px;
  }
  & h3 {
    font-size: 24px;
    font-weight: 500;
    letter-spacing: -0.06px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(6)}px;
  }
  & h4 {
    font-size: 20px;
    font-weight: 500;
    letter-spacing: -0.06px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(4)}px;
  }
  & h5 {
    font-size: 16px;
    font-weight: 500;
    letter-spacing: -0.05px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(2)}px;
  }
  & h6 {
    font-size: 14px;
    font-weight: 500;
    letter-spacing: -0.05px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(2)}px;
  }
  & li {
    font-size: 14px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-left: ${theme.spacing(4)}px;
  }
  & p {
    font-size: 14px;
    margin-bottom: ${theme.spacing(2)}px;
    & > a {
      color: ${theme.palette.secondary.main};
    }
  }
  `}
`;

const renderers = {
  link: (props) => {
    const { href, children, ...other } = props;

    if (!href.startsWith('http')) {
      return (
        <a
          href={href}
          {...other}
        >
          {children}
        </a>
      );
    }

    return (
      <a
        href={href}
        rel="nofollow noreferrer noopener"
        target="_blank"
        {...other}
      >
        {children}
      </a>
    );
  },
  code: (props) => {
    const { language, value, ...other } = props;

    return (
      <SyntaxHighlighter
        language={language}
        style={dracula}
        {...other}
      >
        {value}
      </SyntaxHighlighter>
    );
  }
};

const Docs = () => {
  const history = useHistory();
  const theme = useTheme();
  const { pathname } = useLocation();
  const [file, setFile] = useState(null);

  useEffect(() => {
    const getFile = async () => {
      try {
        // Allow only paths starting with /docs.
        // If you'll use this on another route, remember to check this part.
        if (!pathname.startsWith('/docs')) {
          history.push("/404")
          return;
        }

        const response = await fetch(`/static${pathname}.md`, {
          headers: {
            accept: 'text/markdown' // Do not accept anything else
          }
        });

        if (response.status !== 200) {
          history.push(response.status === 404
            ? '/404'
            : '/500'
          );
          return;
        }

        const data = await response.text();
        setFile(matter(data));
      } catch (err) {
        console.error(err);
        history.push('/500')
      }
    };

    getFile();
  }, [pathname]);

  if (!file) {
    return null;
  }

  return (
    <div>
      <Container
        maxWidth={false}
        style={{ paddingBottom: '120px' }}
      >
        <MarkdownWrapper theme={theme}>
          <Markdown
            escapeHtml
            renderers={renderers}
            source={file.content}
          />
        </MarkdownWrapper>
      </Container>
    </div>
  );
};

export default Docs;
