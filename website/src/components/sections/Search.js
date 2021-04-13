import React, { useState  } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { SectionProps } from '../../utils/SectionProps';
import Input from '../elements/Input';
import DatatablePage from './partials/DatatablePage';

const propTypes = {
  ...SectionProps.types,
  split: PropTypes.bool
}

const defaultProps = {
  ...SectionProps.defaults,
  split: false
}

const Search = ({
  className,
  topOuterDivider,
  bottomOuterDivider,
  topDivider,
  bottomDivider,
  hasBgColor,
  invertColor,
  split,
  ...props
}) => {

const [image, setImage] = useState('')

  const handleTest = (e) => {
    if (e.charCode === 13) {
        setImage(e.target.value)
      }
  }


  const outerClasses = classNames(
    'cta section center-content-mobile reveal-from-bottom',
    topOuterDivider && 'has-top-divider',
    bottomOuterDivider && 'has-bottom-divider',
    hasBgColor && 'has-bg-color',
    invertColor && 'invert-color',
    className
  );


  return (
    <section
      {...props}
      className={outerClasses}
    >
      <div className="container">
        <h1 className="mt-0 mb-16 reveal-from-bottom" data-reveal-delay="200">Search for vulnerabilities</h1>
        <Input onKeyPress={handleTest} type="text" label="Search Image" labelHidden hasIcon="right" placeholder="You Container Image: <domain>/<org>/<image>:<tag>">
            <svg width="16" height="12" xmlns="http://www.w3.org/2000/svg">
            <path d="M9 5H1c-.6 0-1 .4-1 1s.4 1 1 1h8v5l7-6-7-6v5z" fill="#376DF9" />
            </svg>
        </Input>
        <br></br>
        <DatatablePage image={image} />
      </div>
    </section>
  );
}

Search.propTypes = propTypes;
Search.defaultProps = defaultProps;

export default Search;