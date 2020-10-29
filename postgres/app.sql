CREATE OR REPLACE FUNCTION valida_cpf(p_cpf IN CHARACTER VARYING,
                                      p_valida_nulo IN BOOLEAN DEFAULT FALSE)
    RETURNS BOOLEAN AS
$$
DECLARE

    v_cpf_invalidos        CHARACTER VARYING[10]
        DEFAULT ARRAY ['00000000000', '11111111111',
        '22222222222', '33333333333',
        '44444444444', '55555555555',
        '66666666666', '77777777777',
        '88888888888', '99999999999'];
    v_cpf_quebrado         SMALLINT[];
    c_posicao_dv1 CONSTANT SMALLINT DEFAULT 10;
    v_arranjo_dv1          SMALLINT[9] DEFAULT ARRAY [10,9,8,7,6,5,4,3,2];
    v_soma_dv1             SMALLINT DEFAULT 0;
    v_resto_dv1            DOUBLE PRECISION DEFAULT 0;
    c_posicao_dv2 CONSTANT SMALLINT DEFAULT 11;
    v_arranjo_dv2          SMALLINT[10] DEFAULT ARRAY [11,10,9,8,7,6,5,4,3,2];
    v_soma_dv2             SMALLINT DEFAULT 0;
    v_resto_dv2            DOUBLE PRECISION DEFAULT 0;

BEGIN
    IF p_valida_nulo AND nullif(p_cpf, '') IS NULL THEN
        RETURN TRUE;
    END IF;
    IF (NOT (p_cpf ~* '^([0-9]{11})$' OR
             p_cpf ~* '^([0-9]{3}\.[0-9]{3}\.[0-9]{3}\-[0-9]{2})$')
           ) OR
       p_cpf = ANY (v_cpf_invalidos) OR
       p_cpf IS NULL
    THEN
        RETURN FALSE;
    END IF;

    v_cpf_quebrado := regexp_split_to_array(
            regexp_replace(p_cpf, '[^0-9]', '', 'g'), '');
    FOR t IN 1..9
        LOOP
            v_soma_dv1 := v_soma_dv1 +
                          (v_cpf_quebrado[t] * v_arranjo_dv1[t]);
        END LOOP;
    v_resto_dv1 := ((10 * v_soma_dv1) % 11) % 10;

    IF (v_resto_dv1 != v_cpf_quebrado[c_posicao_dv1])
    THEN
        RETURN FALSE;
    END IF;

    FOR t IN 1..10
        LOOP
            v_soma_dv2 := v_soma_dv2 +
                          (v_cpf_quebrado[t] * v_arranjo_dv2[t]);
        END LOOP;
    v_resto_dv2 := ((10 * v_soma_dv2) % 11) % 10;

    RETURN (v_resto_dv2 = v_cpf_quebrado[c_posicao_dv2]);

END;
$$ LANGUAGE plpgsql;

DROP TABLE IF EXISTS consumption;
CREATE TABLE consumption
(
    cpf                  VARCHAR(11) NOT NULL  CHECK (valida_cpf(cpf)),
    private              SMALLINT    NOT NULL,
    incompleto           SMALLINT    NOT NULL,
    data_ultima_compra   DATE  CHECK ( data_ultima_compra <= now() ),
    ticket_medio         DECIMAL  CHECK ( ticket_medio >= 0 ),
    ticket_ultima_compra DECIMAL  CHECK ( ticket_ultima_compra >= 0 ),
    loja_frequente       VARCHAR(14),
    loja_ultima_compra   VARCHAR(14)
);

CREATE OR REPLACE FUNCTION sanitize_cpf_cnpj() RETURNS TRIGGER AS
$$
BEGIN
    NEW.cpf = regexp_replace(NEW.cpf, '[.\-/]', '', 'g');
    NEW.loja_frequente = regexp_replace(NEW.loja_frequente, '[.\-/]', '', 'g');
    NEW.loja_ultima_compra = regexp_replace(NEW.loja_ultima_compra, '[.\-/]', '', 'g');
    RETURN NEW;
END
$$
    LANGUAGE 'plpgsql';

CREATE TRIGGER sanitize_cpf_cnpj_trigger
    BEFORE INSERT
    ON consumption
    FOR EACH ROW
EXECUTE PROCEDURE sanitize_cpf_cnpj()