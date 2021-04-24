import React, { useState, useEffect } from 'react';
import { useLocation, useHistory } from 'react-router-dom';
import Markdown from 'react-markdown';
import matter from 'gray-matter';
import PerfectScrollbar from 'react-perfect-scrollbar';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import dracula from 'react-syntax-highlighter/dist/cjs/styles/prism/dracula';
import { Container, useTheme, makeStyles } from '@material-ui/core';
import styled from 'styled-components';
import Page from 'src/components/Page';

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
    color: #f50057;
    font-family: Inconsolata, Monaco, Consolas, \'Courier New\', Courier monospace;
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
    font-size: 15px;
    font-weight: 500;
    letter-spacing: -0.05px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-top: ${theme.spacing(2)}px;
  }
  & li {
    font-size: 16px;
    margin-bottom: ${theme.spacing(2)}px;
    margin-left: ${theme.spacing(4)}px;
  }
  & p {
    font-size: 16px;
    margin-bottom: ${theme.spacing(2)}px;
    & > a {
      color: ${theme.palette.secondary.main};
    }
  }
  `}
`;

const renderers = {
  // root: ({ children }) => {
  //   const TOCLines = children.reduce((acc, { key, props }) => {
  //     // Skip non-headings
  //     if (key.indexOf('heading') !== 0) {
  //       return acc;
  //     }

  //     // Indent by two spaces per heading level after h1
  //     let indent = '';
  //     for (let idx = 1; idx < props.level; idx++) {
  //       indent = `${indent}  `;
  //     }

  //     // Append line to TOC
  //     // This is where you'd add a link using Markdown syntax if you wanted
  //     return acc.concat([`${indent}* ${props.node.children[0].value}`]);
  //   }, []);

  //   console.log(TOCLines)
  //   return (
  //     <div>
  //       <Markdown source={TOCLines.join("\n")} />
  //       {children}
  //     </div>
  //   );
  // },
  link: (props) => {
    const { href, children, ...other } = props;

    if (!href.startsWith('http')) {
      return (
        <a
          href={href}
          target="_blank"
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
  },
  heading: (props) => {
    function flatten(text, child) {
      return typeof child === 'string'
        ? text + child
        : React.Children.toArray(child.props.children).reduce(flatten, text)
    }
    var children = React.Children.toArray(props.children)
    var text = children.reduce(flatten, '')
    var slug = text.toLowerCase().replace(/\W/g, '-')
    return React.createElement('h' + props.level, { id: slug }, props.children)
  },
  image: ({
    alt,
    src,
    title
  }) => (
    <img
      alt={alt}
      src={src}
      title={title}
      style={{ maxWidth: "100%" }} />
  )
};

const Docs = () => {
  const history = useHistory();
  const theme = useTheme();
  const { pathname } = useLocation();
  const [file, setFile] = useState(null);


  useEffect(() => {
    document.getElementById("content-container").scrollTo(0, 0);

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
      <Page title={`${file.data.title}`}>

        <Container
          maxWidth={"lg"}
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
      </Page>
    </div>
  );
};

export default Docs;
